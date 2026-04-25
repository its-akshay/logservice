package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("i am processing now id: %d job %d \n", id, job)
	}
}

func producer(jobs chan<- int) {
	for i := 1; i <= 100; i++ {
		jobs <- i
	}
	close(jobs) 
}

func main() {
	jobs := make(chan int, 20)
	var wg sync.WaitGroup

	numWorkers := 5
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	//sending jobs
	go producer(jobs)
	
	wg.Wait()
}
