package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	serverIP    = "127.0.0.1"
	serverPort  = "9000"
	numRequests = 100000
	numWorkers  = 100
	command     = "instruction=get|store=user_data|key=5"
)

func main() {
	results := make(chan float64, numRequests)
	startBenchmark := time.Now()

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", serverIP+":"+serverPort)
			if err != nil {
				fmt.Println("Connection error:", err)
				return
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)

			for i := 0; i < numRequests/numWorkers; i++ {
				start := time.Now()

				_, err := fmt.Fprintf(conn, "%s\n", command)
				if err != nil {
					fmt.Println("Send error:", err)
					return
				}

				_, err = reader.ReadString('\n')
				if err != nil {
					fmt.Println("Read error:", err)
					return
				}

				results <- time.Since(start).Seconds()
			}
		}()
	}

	wg.Wait()
	close(results)

	var latencies []float64
	for l := range results {
		latencies = append(latencies, l)
	}

	if len(latencies) == 0 {
		fmt.Println("No successful requests, exiting benchmark.")
		return
	}

	total := 0.0
	maxLatency := 0.0
	minLatency := 1e9
	for _, l := range latencies {
		total += l
		if l > maxLatency {
			maxLatency = l
		}
		if l < minLatency {
			minLatency = l
		}
	}

	elapsed := time.Since(startBenchmark).Seconds()

	fmt.Println("\n=== Benchmark Results ===")
	fmt.Printf("Total requests:   %d\n", len(latencies))
	fmt.Printf("Workers:          %d\n", numWorkers)
	fmt.Printf("Average latency:  %.6f sec\n", total/float64(len(latencies)))
	fmt.Printf("Min latency:      %.6f sec\n", minLatency)
	fmt.Printf("Max latency:      %.6f sec\n", maxLatency)
	fmt.Printf("Throughput:       %.2f req/sec\n", float64(len(latencies))/elapsed)
	fmt.Printf("Total time:       %.6f sec\n", elapsed)
}
