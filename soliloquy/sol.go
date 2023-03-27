// soliloquy —from the Latin solus ("alone") and loqui ("to speak")—
// is a speech that one gives to oneself.

// This program launch a web server and a web client, thus speaking
// to itself

package main

import (
	"fmt"
	"io"
	"log"
	_ "net"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/google/uuid"
	_ "github.com/jsimonetti/rtnetlink/rtnl"
	_ "github.com/milosgajdos/tenus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func linkSetUp( space netns.NsHandle, link netlink.Link, index, host int) {
	nsOrig, err := netns.Get()
	if err != nil {
		log.Fatal("Can't get original netspace: ", err)
	}
	defer netns.Set(nsOrig)

	if err := netns.Set(space); err != nil {
		log.Fatal("Can't switch to namespace: ")
	}

	addr, _ := netlink.ParseAddr(fmt.Sprintf("10.0.%d.%d/24", index, host))
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		log.Fatal("Can't set IP address: ", err)
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		log.Fatal("Can't bring up client address: ", err)
	}
}

func createVeth(clientSpace netns.NsHandle, clientIndex int, simSpace netns.NsHandle) {
	nsOrig, err := netns.Get()
	if err != nil {
		log.Fatal("Can't get original netspace: ", err)
	}
	defer netns.Set(nsOrig)


	err = netns.Set(simSpace)
	if err != nil {
		log.Fatal("Can't get original netspace: ", err)
	}

	// Thanks, https://github.com/teddyking/netsetgo/blob/0.0.1/device/veth.go#L16

	vethLinkAttrs := netlink.NewLinkAttrs()
	vethLinkAttrs.Name = fmt.Sprintf("client%d", clientIndex)

	veth := &netlink.Veth{
		LinkAttrs: vethLinkAttrs,
		PeerName: "sim",
	}

	if err := netlink.LinkAdd(veth); err != nil {
		log.Fatal("Can't create veth pair: ", err)
	}

	// Put the endpoints in the correct spaces
	simLink, _ := netlink.LinkByName(veth.Name)
	clientLink, _ := netlink.LinkByName(veth.PeerName)

	// netns.NsHandle is just a file descriptor in disguise
	err = netlink.LinkSetNsFd(clientLink, int(clientSpace))
	if err != nil {
		log.Fatal("Can't put pair into client namespace: ", err)
	}

	linkSetUp(clientSpace, clientLink, clientIndex, 2)
	linkSetUp(simSpace, simLink, clientIndex, 1)

	// All done!
}

func createNSWithLoopback(name, suffix string) netns.NsHandle {
	nsOrig, err := netns.Get()
	if err != nil {
		log.Fatal("Can't get original netspace: ", err)
	}
	defer netns.Set(nsOrig)

	// NewNamed leaves us in the new namespace
	nsHandle, err := netns.NewNamed(name + suffix)
	if err != nil {
		log.Fatal("Can't create namespace: ", err)
	}

	lo, err := netlink.LinkByName("lo")
	if err != nil {
		log.Fatal("Can't find loopback: ", err)
	}

	err = netlink.LinkSetUp(lo)
	if err != nil {
		log.Fatal("Can't find loopback: ", err)
	}

	return nsHandle
}

// Create all of the required namespaces,
// veth pairs, and bridges. Assign
// addresses as appropriate.
func config() {
	// TODO completely refactor and generalize

	// So we can distinguish test runs
	nsSuffix := "-" + uuid.New().String()

	var clientSpaces []netns.NsHandle

	clientSpace := createNSWithLoopback("client0", nsSuffix)
	clientSpaces = append(clientSpaces, clientSpace)

	clientSpace = createNSWithLoopback("client1", nsSuffix)
	clientSpaces = append(clientSpaces, clientSpace)

	simSpace := createNSWithLoopback("sim", nsSuffix)

	// For each client, create a veth pair between
	// it and the sim
	for i, clientSpace := range(clientSpaces) {
		createVeth(clientSpace, i, simSpace)
	}
}


// Run an http server in the background, then
// run a set of http clients in the foreground.
// The server prints the message that the client sent.
// The client prints the message that the server sent.
func main() {
	config()
	return
	fmt.Println("yo, wo")
	go server()

	// TODO One really shouldn't use `Sleep` to manage
	// concurrency.
	time.Sleep(time.Second * 1)

	u := url.URL{
		Host:   "localhost:8000",
		Scheme: "http",
		Path:   "/",
	}
	for i := 0; i < 500; i++ {
		u.RawQuery = url.Values{
			"i":   {fmt.Sprint(i)},
			"msg": {"Hello from client"},
		}.Encode()
		clientInternal(u)

		u.RawQuery = url.Values{
			"i":   {fmt.Sprint(i * 100)},
			"msg": {"Hello from cURL"},
		}.Encode()
		clientExternal(u)
	}
}

// Fetch the named `URL` and display the returned `.Body`
func clientInternal(u url.URL) {
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("C Response failed: %d %s", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("I %s\n", body)
}

// Fetch the named `URL` using the external `curl` command
func clientExternal(u url.URL) {
	out, err := exec.Command("curl", u.String()).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("X %s\n", out)
}

// Launch a web server on "http://0.0.0.0:8000/". For each request,
// decode the parameters and return a message in the response.
func server() {
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			i := r.FormValue("i")
			resp := fmt.Sprintf("%s: %s", i, "Hello from server")
			io.WriteString(w, resp)
			fmt.Printf("S %s: %s\n", i, r.FormValue("msg"))
		})
	http.ListenAndServe(":8000", nil)
}
