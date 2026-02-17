package main

import "fmt"

// struct
type Elevator struct {
	currentFloor int
}

// methon on pointer

// Here we created a method which will get memory add for struct as pointer and then accessing the struct variables aka inheritence
func (e *Elevator) MoveUp() {
	e.currentFloor++
}

func (e *Elevator) MoveDown() {
	e.currentFloor--
}

func main() {

	e := Elevator{currentFloor: 1}
	fmt.Println("Now at floor:", e.currentFloor)

	e.MoveUp()
	fmt.Println("Now at floor:", e.currentFloor)

	e.MoveUp()
	fmt.Println("Now at floor:", e.currentFloor)

	e.MoveDown()
	fmt.Println("Now at floor:", e.currentFloor)

}
