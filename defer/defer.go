package main

import "fmt"
import "math/rand"
import "time"

func main() {
	rand.Seed(time.Now().UnixNano())
	messages := []string{ "a", "b", "c" }
	fmt.Println("Hello!")
	for i := 0; i < 5; i++ {
		message := messages[rand.Intn(len(messages))]
		fmt.Printf("%d: %v\n", i, message)
		if i % 2 == 0 {
			defer fmt.Printf("%d: %v\n", i, message)
		}
	}
	fmt.Println("Goodbye!")
}
