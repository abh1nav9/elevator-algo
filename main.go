package main

import (
	"fmt"
	"time"
)

// struct
type Elevator struct {
	currentFloor int
	requests     []int
}

func (e *Elevator) step() {
	if len(e.requests) == 0 {
		return
	}

	target := e.requests[0]

	if e.currentFloor < target {
		e.currentFloor++
	} else if e.currentFloor > target {
		e.currentFloor--
	} else {
		fmt.Println("Reached floor:", target)
		// e.requests[1:] is used to remove first element from []requests
		e.requests = e.requests[1:]
	}
}

func main() {

	e := Elevator{currentFloor: 1, requests: []int{4, 10, 3}}

	for {
		if len(e.requests) == 0 {
			fmt.Println("All requests completed. Elevator is idle at floor:", e.currentFloor)
			break
		}

		e.step()
		fmt.Println("Current Floor:", e.currentFloor, "| Requests:", e.requests)
		time.Sleep(800 * time.Millisecond)
	}
}
