package main

/* This code demonstrates the use of a mutex to safely update a shared variable
   from multiple goroutines. The `updateMessage` function locks the mutex before
   updating the shared variable `msg`, ensuring that only one goroutine can modify
   `msg` at a time. This is crucial to avoid the race condition that can occur
   when multiple goroutines attempt to read and write to the same variable concurrently.

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
*/

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
