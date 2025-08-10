package main

import (
	"fmt"
)

func printSomething(s string) {
	fmt.Println(s)
}

func main() {
	go printSomething("This is the first thing to be printed!")

	/* bad practice: using time.Sleep to wait for goroutines */
	// time.Sleep(1 * time.Second)

	printSomething("This is the second thing to be printed!")
}
