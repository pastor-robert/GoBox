// rsg0 is a discrete event simulation inspired by SystemC.
// It provides events, notifications, time, channels, etc.
package main

import (
	"fmt"
	"time"
)

// generic code follows

// RSG is a descrete-event simulator inspired by SystemC
type RSG struct {
	Name string
}

// Start a simulation and let it run until
// there are no notificaitons pending.
// Returns after the simulation is complete.
func (s *RSG) Start() {
}

// Start a simulation and let it run until
// there are no notifications pending,
// or until the indicated simulated nanoseconds pass,
// whichever comes first.
func (s *RSG) StartTimeout(d time.Duration) {
}


// The Event primitive allows one to wait:
//   for a specific amount of time
//   for a specific notification
//   for a specific notification, with timeout
type Event struct {
	sim *RSG
	Name string
}

// Wait for this event until the cows come home
func (e Event) Wait() {
}

// Wait for this event until the indicated time
func (e Event) WaitTimeout(t int) {
}



// App-specific code follows:
type HelloWorld struct {
	RSG
	ev *Event
}

func NewHelloWorld() *HelloWorld {
	h := &HelloWorld{}
	h.ev = &Event{sim: &h.RSG, Name: "hello"}
	return h
}



func main() {
	// sim := &RSG{Name: "My Simulator"}
	hello := NewHelloWorld()
	hello.Start()
	fmt.Println("Goodbye")
}
