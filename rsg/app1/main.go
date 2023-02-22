// RSG ("Ready, Sim, Go") is a simulation engine inspired by SystemC

package main

import (
	"fmt"
	"rob-adams.us/rsg"
)

func hey(s *rsg.Simulation, e *rsg.Event, x any) {
	i := x.(int)
	fmt.Printf("%d\n", i)
	if i < 1e7 {
		s.Wait(e, hey, i*10)
	}
}

func main() {
	s := rsg.NewSimulation("My Sim")
	e1 := s.NewEvent("EV1")
	e2 := s.NewEvent("EV2")
	s.Wait(e1, hey, 1)
	s.Wait(e2, hey, 2)
	fmt.Printf("------------\n")
	s.DeltaOne()
	fmt.Printf("------------\n")
	s.DeltaOne()
	fmt.Printf("------------\n")
	s.DeltaOne()
	fmt.Printf("------------\n")
	s.DeltaN()
}
