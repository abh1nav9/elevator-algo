package main

import "fmt"

type Elevator struct {
	currentFloor int
}

func main() {

	e := Elevator{
		currentFloor: 1,
	}

	fmt.Println("Elevator is at floor:", e.currentFloor)
}