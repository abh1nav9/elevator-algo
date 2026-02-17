package main

import "fmt"

func main() {
	floors := 10

	for i := floors; i >= 1; i-- {
		fmt.Println("Floor", i)
	}
}
