package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type direction int

const (
	IDLE direction = iota
	UP
	DOWN
)

// struct
type Elevator struct {
	mu           sync.Mutex
	currentFloor int
	requests     []int
	direction    direction
}

func (e *Elevator) addRequest(floor int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, f := range e.requests {
		if f == floor {
			return
		}
	}

	e.requests = append(e.requests, floor)
}

func (e *Elevator) step() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.requests) == 0 {
		e.direction = IDLE
		return
	}

	target := e.requests[0]

	if e.currentFloor < target {
		e.direction = UP
		e.currentFloor++
	} else if e.currentFloor > target {
		e.direction = DOWN
		e.currentFloor--
	} else {
		fmt.Println("Reached floor:", target)
		e.requests = e.requests[1:]
		if len(e.requests) == 0 {
			e.direction = IDLE
		}
	}
}

func (d direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	default:
		return "IDLE"
	}
}

func render(e *Elevator, maxFloors int) {
	e.mu.Lock()
	currentFloor := e.currentFloor
	queue := append([]int{}, e.requests...)
	dir := e.direction
	e.mu.Unlock()

	// Only clear screen if elevator is moving
	if len(queue) > 0 {
		fmt.Print("\033[H\033[2J")
	}

	for i := maxFloors; i >= 1; i-- {
		if i == currentFloor {
			fmt.Printf("Floor %2d | [ E ]\n", i)
		} else {
			fmt.Printf("Floor %2d | [   ]\n", i)
		}
	}

	fmt.Println()

	if len(queue) == 0 {
		fmt.Println("Status   : IDLE")
		fmt.Println("Direction:", dir)
		fmt.Println("Type     : <floor> or go <floor> (e.g., 4 or go 4)")
	} else {
		fmt.Println("Queue    :", queue)
		fmt.Println("Direction:", dir)
	}
}

func main() {

	e := Elevator{currentFloor: 1, direction: IDLE}
	maxFloors := 10
	lastFloor := -1
	lastQueueLen := -1

	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" {
				fmt.Println("Bye 👋")
				os.Exit(0)
			}

			// Allow just a number: "4"
			if floor, err := strconv.Atoi(input); err == nil {
				if floor >= 1 && floor <= maxFloors {
					e.addRequest(floor)
					continue
				}
				fmt.Println("Invalid floor (1-10 only)")
				continue
			}

			// Allow: "go 4"
			parts := strings.Fields(input)
			if len(parts) == 2 && parts[0] == "go" {
				floor, err := strconv.Atoi(parts[1])
				if err == nil && floor >= 1 && floor <= maxFloors {
					e.addRequest(floor)
				} else {
					fmt.Println("Invalid floor (1-10 only)")
				}
				continue
			}

			fmt.Println("Use: <floor> or go <floor> (e.g., 4 or go 4)")
		}
	}()

	for {
		e.step()

		e.mu.Lock()
		changed := e.currentFloor != lastFloor || len(e.requests) != lastQueueLen
		lastFloor = e.currentFloor
		lastQueueLen = len(e.requests)
		e.mu.Unlock()

		if changed {
			render(&e, maxFloors)
		}

		time.Sleep(700 * time.Millisecond)
	}
}
