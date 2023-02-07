package greetings

import (
	"fmt"
	"errors"
)

// Return a greeting for the named person
func Hello(name string) (string, error) {
	if name == "" {
		return "", errors.New("name not specified")
	}
	return fmt.Sprintf("Hi, %v. Welcome!", name), nil
}
