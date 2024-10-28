package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	url         string
	requests    int
	concurrency int
)

func init() {
	flag.StringVar(&url, "url", "", "Service URL")
	flag.IntVar(&requests, "requests", 100, "Total number of requests")
	flag.IntVar(&concurrency, "concurrency", 10, "Number of concurrent requests")
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Println("Service URL is required.")
		return
	}

	runStressTest(url, requests, concurrency)
}

func runStressTest(url string, requests, concurrency int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalRequests int
	var status200 int
	statusDistribution := make(map[int]int)

	start := time.Now()

	ch := make(chan struct{}, concurrency)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		ch <- struct{}{}

		go func(requestNumber int) {
			defer wg.Done()
			resp, err := http.Get(url)

			<-ch

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fmt.Printf("Error on request #%d: %v\n", requestNumber+1, err)
				return
			}
			defer resp.Body.Close()

			totalRequests++
			statusDistribution[resp.StatusCode]++

			if resp.StatusCode == 200 {
				status200++
			}

			if totalRequests%100 == 0 {
				fmt.Printf("Progress: %d requests completed\n", totalRequests)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("Load Test Report:")
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Status 200: %d\n", status200)
	fmt.Println("HTTP Status Distribution:")

	for status, count := range statusDistribution {
		fmt.Printf("Status %d: %d\n", status, count)
	}
}
