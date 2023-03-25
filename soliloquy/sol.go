package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Run an http server in the background, then
// run a series of http clients in the foreground.
// The server prints the message that the client sent.
// The client prints the message that the server sent.
func main() {
	fmt.Println("yo, wo")
	go server()

	// TODO One really shouldn't use `Sleep` to manage
	// concurrency.
	time.Sleep(time.Second * 1)

	// TODO loop count should be configurable

	for i := 0; i < 500; i++ {
		u := url.URL{
			Host:   "localhost:8000",
			Scheme: "http",
			Path:   "/",
			RawQuery: url.Values{
				"msg": {"Hello from client"},
				"i":   {fmt.Sprint(i)},
			}.Encode(),
		}
		client(u)
	}
}

func client(u url.URL) {
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
	fmt.Printf("C %s\n", body)
}

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
