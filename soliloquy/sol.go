// soliloquy —from the Latin solus ("alone") and loqui ("to speak")—
// is a speech that one gives to oneself.

// This program launch a web server and a web client, thus speaking
// to itself

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"time"
)

// Run an http server in the background, then
// run a set of http clients in the foreground.
// The server prints the message that the client sent.
// The client prints the message that the server sent.
func main() {
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
			"i": {fmt.Sprint(i)},
			"msg": {"Hello from client"},
		}.Encode()
		clientInternal(u)

		u.RawQuery = url.Values{
			"i": {fmt.Sprint(i*100)},
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
