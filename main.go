package main

import "fmt"

// struct
type Elevator struct {
	currentFloor int
	requests     []int
}

// methon on pointer
func (e *Elevator) addRequests(floor int) {
	// append - add floor to requests array
	e.requests = append(e.requests, floor)
}

func main() {

	e := Elevator{currentFloor: 1}

	e.addRequests(4)
	e.addRequests(2)
	e.addRequests(8)

	fmt.Println("Requests:", e.requests)
}
