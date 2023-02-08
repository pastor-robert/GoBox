package main

import "fmt"
import "math/rand"

func main() {
	messages := []string{ "a", "b", "c" }
	fmt.Println("Hello!")
	for i := 1; i <= 5; i++ {
		message := messages[rand.Intn(len(messages))]
		fmt.Printf("%d: %v\n", i, message)
		if i % 2 == 0 {
			defer fmt.Printf("%d: %v\n", i, message)
		}
	}
	fmt.Println("Goodbye!")
}
