package main

import (
	"fmt"
	"log"
	"github.com/pastor-robert/GoBox/greetings"
)


func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	message, err := greetings.Hello("Ruth")
	fmt.Println(message)

	message, err = greetings.Hello("")
	log.Fatal(err)
}
