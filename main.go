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

func render(e *Elevator, maxFloors int) {
	fmt.Print("\033[H\033[2J")

	for i := maxFloors; i >= 1; i-- {
		if i == e.currentFloor {
			fmt.Printf("Floor %2d | [ E ]\n", i)
		} else {
			fmt.Printf("Floor %2d | [   ]\n", i)
		}
	}

	fmt.Println("\nQueue:", e.requests)
}

func main() {

	e := Elevator{currentFloor: 1, requests: []int{4, 10, 3}}

	for {
		if len(e.requests) == 0 {
			fmt.Println("All requests completed. Elevator is idle at floor:", e.currentFloor)
			break
		}

		e.step()
		render(&e, 10)
		time.Sleep(700 * time.Millisecond)
	}
}
