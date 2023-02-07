package greetings

import "fmt"

// Return a greeting for the named person
func Hello(name string) string {
	return fmt.Sprintf("Hi, %v. Welcome!", name)
}
