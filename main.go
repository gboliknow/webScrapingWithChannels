package main

import (
	"fmt"
	"net/http"
	"time"
)

// Task represents a unit of work that needs to be done.
type Task struct {
	url string
}

// Result represents the result of a task.
type Result struct {
	url      string
	response *http.Response
	err      error
}

// worker is a function that processes tasks from the tasks channel and sends results to the results channel.
func worker(id int, tasks <-chan Task, results chan<- Result, timing chan<- time.Duration) {
	for task := range tasks {
		start := time.Now()
		fmt.Printf("Worker %d started task for %s\n", id, task.url)
		resp, err := http.Get(task.url)
		results <- Result{url: task.url, response: resp, err: err}
		elapsed := time.Since(start)
		timing <- elapsed
		fmt.Printf("Worker %d finished task for %s in %v\n", id, task.url, elapsed)
	}
}

func main() {
	// Define the list of URLs to be scraped.
	urls := []string{
		"http://example.com",
		"http://google.com",
		"http://golang.org",
		"http://github.com",
	}

	// Create channels for tasks and results.
	tasks := make(chan Task, len(urls))
	results := make(chan Result, len(urls))
	timing := make(chan time.Duration, len(urls))
	// done := make(chan struct{})

	// Start a pool of workers.
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, results, timing)
	}

	// Send tasks to the task channel.
	for _, url := range urls {
		tasks <- Task{url: url}
	}
	close(tasks)

	// Collect results from the results channel.
	for i := 0; i < len(urls); i++ {
		result := <-results
		if result.err != nil {
			fmt.Printf("Error fetching %s: %v\n", result.url, result.err)
		} else {
			fmt.Printf("Fetched %s with status %s\n", result.url, result.response.Status)
		}
	}

	// Close the results channel after all results have been processed.
	close(results)

	// Give some time for all goroutines to finish before the program exits.
	time.Sleep(2 * time.Second)

	// //using select to process result

	// go func() {
	// 	for i := 0; i < len(urls); i++ {
	// 		result := <-results
	// 		if result.err != nil {
	// 			fmt.Printf("Error fetching %s: %v\n", result.url, result.err)
	// 		} else {
	// 			fmt.Printf("Fetched %s with status %s\n", result.url, result.response.Status)
	// 		}
	// 	}
	// 	done <- struct{}{}
	// }()

	// select {
	// case <-done:
	// 	fmt.Println("All tasks are done")
	// case <-time.After(5 * time.Second):
	// 	fmt.Println("Timeout waiting for tasks to complete")
	// }
}
