package main

import (
	"fmt"
	"sync"
)

var msg string

func updateMessage(s string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	msg = s
}

func printMessage() {
	fmt.Println(msg)
}

func main() {
	var wg sync.WaitGroup

	msg = "Hello, world!"

	go updateMessage("Hello, universe!", &wg)
	printMessage()

	go updateMessage("Hello, cosmos!", &wg)
	printMessage()

	go updateMessage("Hello, world!", &wg)
	printMessage()
}
