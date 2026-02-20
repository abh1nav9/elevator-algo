package main

import (
	"fmt"
	"sort"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type direction int

const (
	IDLE direction = iota
	UP
	DOWN
)

type Elevator struct {
	mu           sync.Mutex
	currentFloor int
	requests     []int
	direction    direction
	minFloor     int
	maxFloor     int
}

func NewElevator(minFloor, maxFloor int) *Elevator {
	return &Elevator{
		currentFloor: minFloor,
		direction:    IDLE,
		minFloor:     minFloor,
		maxFloor:     maxFloor,
	}
}

func (e *Elevator) AddRequest(floor int) bool {
	if floor < e.minFloor || floor > e.maxFloor {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	for _, f := range e.requests {
		if f == floor {
			return false
		}
	}

	if floor == e.currentFloor && e.direction == IDLE {
		return false
	}

	e.requests = append(e.requests, floor)
	return true
}

func (e *Elevator) nextTarget() (int, bool) {
	if len(e.requests) == 0 {
		return 0, false
	}

	curr := e.currentFloor
	var above, below []int

	for _, r := range e.requests {
		if r == curr {
			return r, true
		}
		if r > curr {
			above = append(above, r)
		} else {
			below = append(below, r)
		}
	}

	sort.Ints(above)
	sort.Sort(sort.Reverse(sort.IntSlice(below)))

	switch e.direction {
	case UP:
		if len(above) > 0 {
			return above[0], true
		}
		if len(below) > 0 {
			return below[0], true
		}

	case DOWN:
		if len(below) > 0 {
			return below[0], true
		}
		if len(above) > 0 {
			return above[0], true
		}

	case IDLE:
		if len(above) == 0 && len(below) == 0 {
			return 0, false
		}
		var closest int
		var minDist int

		if len(above) > 0 {
			closest = above[0]
			minDist = above[0] - curr
		}
		if len(below) > 0 {
			dist := curr - below[0]
			if len(above) == 0 || dist < minDist {
				closest = below[0]
				minDist = dist
			}
		}

		return closest, true
	}

	return 0, false
}

func (e *Elevator) Step() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.requests) == 0 {
		e.direction = IDLE
		return
	}

	target, ok := e.nextTarget()
	if !ok {
		e.direction = IDLE
		return
	}

	if target > e.currentFloor {
		e.direction = UP
		e.currentFloor++
	} else if target < e.currentFloor {
		e.direction = DOWN
		e.currentFloor--
	} else {
		e.removeRequest(target)
		if len(e.requests) == 0 {
			e.direction = IDLE
		}
	}
}

func (e *Elevator) removeRequest(floor int) {
	newReqs := e.requests[:0]
	for _, r := range e.requests {
		if r != floor {
			newReqs = append(newReqs, r)
		}
	}
	e.requests = newReqs
}

func (e *Elevator) State() (int, []int, direction) {
	e.mu.Lock()
	defer e.mu.Unlock()
	reqs := make([]int, len(e.requests))
	copy(reqs, e.requests)
	return e.currentFloor, reqs, e.direction
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

type tickMsg time.Time

type model struct {
	elevator  *Elevator
	maxFloors int
}

func initialModel(elevator *Elevator, maxFloors int) model {
	return model{
		elevator:  elevator,
		maxFloors: maxFloors,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.elevator.Step()
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			floor := int(msg.String()[0] - '0')
			if floor <= m.maxFloors {
				m.elevator.AddRequest(floor)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	current, queue, dir := m.elevator.State()

	shaft := ""
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(7).
		Align(lipgloss.Center)

	for i := m.maxFloors; i >= 1; i-- {
		floorLabel := fmt.Sprintf("Floor %2d", i)
		if i == current {
			shaft += fmt.Sprintf("%s %s\n", floorLabel, boxStyle.Render(" E "))
		} else {
			shaft += fmt.Sprintf("%s %s\n", floorLabel, boxStyle.Render("   "))
		}
	}

	info := lipgloss.NewStyle().Bold(true).Render("Elevator TUI\n\n")
	info += fmt.Sprintf("Current: Floor %d\n", current)
	info += fmt.Sprintf("Direction: %s\n", dir)
	info += fmt.Sprintf("Queue: %v\n\n", queue)
	info += fmt.Sprintf("Press 1-%d to request floor | q to quit\n", m.maxFloors)

	return lipgloss.JoinHorizontal(lipgloss.Top, shaft, "  ", info)
}

func main() {
	const maxFloors = 9
	elevator := NewElevator(1, maxFloors)

	p := tea.NewProgram(initialModel(elevator, maxFloors), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
