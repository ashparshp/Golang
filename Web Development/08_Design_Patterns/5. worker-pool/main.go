package main

import (
	"fmt"
	"time"
)

// worker is the function that processes a unit of work. It is called
// from the main function as a goroutine, so all workers run concurrently.
func worker(id int, jobs chan int, results chan int) {
	for j := range jobs {
		fmt.Println("Worker", id, "started job", j, "...")
		time.Sleep(time.Second)
		fmt.Println("Worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	// numJobs is the number of jobs we have to perform.
	const numJobs = 5

	// Create two channels: one to send work to, and one
	// that will send results back to us.
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Create our worker pool, with three workers running as goroutines.
	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}

	// Send jobs to our worker pool.
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Print out the results as they come in.
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
