// RSG ("Ready, Sim, Go") is a simulation engine inspired by SystemC

package rsg

import (
	"container/heap"
	"testing"
	"fmt"
)

// import (
// 	"testing"
// )

type SimTime uint64

// An Event is something that can be waited upon.
// An Event can have several waiters, but only one
// deadline.
type Event struct {
	Name     string
	deadline SimTime
}

// A Waiter is a something that waits
// on an Event. The specified method is
// invoked when the deadline arrives
type Waiter struct {
	callback Callback
	x        any
}
type Callback func(*Event, any)

// An EventList associates Waiters with the
// corresponding Event
type WaitList []Waiter
type EventList map[*Event]WaitList

func NewWaitList() WaitList { return WaitList{} }

// A Simulation is a collection of Modules
// which are collections of Waiters, waiting
// on Events.
// The Simulation controls the flow of time.
type Simulation struct {
	Name string

	// The Simulation is responsible for the passage of time.
	now SimTime

	// All Events are stored as keys
	// in the EventList map. Waitlists are
	// the values. When a specific event
	// fires, the Simulator invokes all of the
	// Waiters and clears the wait list
	events EventList

	// All runnable Events are stored in
	// a simple slice
	runnable eventHeap
}

func NewSimulation(n string) *Simulation {
	s := &Simulation{}
	s.Name = n
	s.now = 0
	s.events = EventList{}
	s.runnable = eventHeap{}
	heap.Init(&s.runnable)
	return s
}

// This data structure implements heap.Interface. The data contained
// within is redundant with, and a subset of, s.EventList,
// but this is ordered by e.deadline
type eventHeap []*Event

func (h eventHeap) Len() int           { return len(h) }
func (h eventHeap) Less(i, j int) bool { return h[i].deadline < h[j].deadline }
func (h eventHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *eventHeap) Push(x any) {
	*h = append(*h, x.(*Event))
}
func (h *eventHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Query the Simulation for the current time
func (s *Simulation) Now() SimTime {
	return s.now
}

// Add to an Event's wait list
func (s *Simulation) Wait(e *Event, f Callback, x any) {
	s.events[e] = append(s.events[e], Waiter{f, x})

	// TODO Everything Everywhere All At Once
	e.deadline = s.Now()
	s.EventPush(e)
}

// NewEvent receives a *Simulation so that
// Simulation can update all of the bookkeeping
func (s *Simulation) NewEvent(name string) *Event {
	e := &Event{}
	e.Name = name
	s.events[e] = NewWaitList()
	return e
}

func (s *Simulation) EventPeek() *Event {
	return s.runnable[0]
}

func (s *Simulation) EventPop() *Event {
	return heap.Pop(&s.runnable).(*Event)
}

func (s *Simulation) EventPush(e *Event) {
	heap.Push(&s.runnable, e)
}

// Perform all of the events that expire in the current delta cycle
func (s *Simulation) deltaOne() {

	// Grab all candidates to avoid crazy loops later
	immediate := []*Event{}
	for len(s.runnable) > 0 && s.EventPeek().deadline == s.now {
		e := s.EventPop()
		immediate = append(immediate, e)
	}

	// For each qualifying Event, call all of the
	// waiting functions.
	for _, e := range immediate {

		// Make a copy of the waitlist to avoid
		// crazy loops.
		// TODO: This seems expensive. Use a pointer for WaitList
		wl := s.events[e]
		s.events[e] = NewWaitList()

		// For each Waiter on the Waitlist, call the Waiter
		for _, w := range wl {
			w.callback(e, w.x)
		}
	}
}

func TestXxx(t *testing.T) {
	fmt.Printf("in the test\n")
	t.Errorf("wah-wah")
}
