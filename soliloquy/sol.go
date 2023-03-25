package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("yo, wo")
	go server()
	for i := 0; i < 5; i++ {
		client("http://localhost:8000")
	}
}

func client(url string) {
	resp, err := http.Get(url + "?Hello+from+client")
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
			io.WriteString(w, "Hello from server")
			fmt.Printf("S %s\n", r.URL)
		})
	http.ListenAndServe(":8000", nil)
}
