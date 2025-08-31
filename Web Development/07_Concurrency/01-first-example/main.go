package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure that Done is called when the goroutine finishes
	fmt.Println(s)
}

func main() {

	/* bad practice: using time.Sleep to wait for goroutines
	go printSomething("This is the first thing to be printed!")
	time.Sleep(1 * time.Second)
	*/

	var wg sync.WaitGroup
	words := []string{"alpha!", "beta!", "delta", "gamma!", "pi", "zeta!", "eta", "sigma!", "epsilon!"}
	wg.Add(len(words)) // Set the number of goroutines to wait for

	for i, word := range words {
		go printSomething(fmt.Sprintf("Word %d: %s", i+1, word), &wg)
	}
	wg.Wait() // Wait for all goroutines to finish

	wg.Add(1)
	printSomething("This is the second thing to be printed!", &wg)
}
