package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(newMsg string, mu *sync.Mutex) {
	defer wg.Done()

	mu.Lock()
	msg = newMsg
	mu.Unlock()
}

func main() {
	msg = "Hello, World!"

	var mu sync.Mutex

	wg.Add(2)
	go updateMessage("Hello, Go!", &mu)
	go updateMessage("Hello, Concurrency!", &mu)
	wg.Wait()

	fmt.Println(msg)
}

/* This may not work as expected due to race conditions.
   The variable 'msg' is being updated by multiple goroutines
   without synchronization, which can lead to unpredictable results.
   To fix this, we can use a mutex to ensure that only one goroutine
   can update 'msg' at a time.

var msg string
var wg sync.WaitGroup

func updateMessage(newMsg string) {
	defer wg.Done()
	msg = newMsg
}

func main() {
	msg = "Hello, World!"

	wg.Add(2)
	go updateMessage("Hello, Go!")
	go updateMessage("Hello, Concurrency!")
	wg.Wait()

	fmt.Println(msg)
}
*/
