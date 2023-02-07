package main

import (
	"fmt"
	"github.com/pastor-robert/GoBox/greetings"
)


func main() {
	message := greetings.Hello("Ruth")
	fmt.Println(message)
}
