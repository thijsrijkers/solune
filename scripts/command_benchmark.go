package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	serverIP    = "127.0.0.1"
	serverPort  = "9000"
	numRequests = 1
	command     = ""
)

func main() {
	var latencies []float64

	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	startBenchmark := time.Now()
	for i := 0; i < numRequests; i++ {
		start := time.Now()

		_, err := fmt.Fprintf(conn, "%s\n", command)
		if err != nil {
			fmt.Println("Send error:", err)
			break
		}

		_, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		latencies = append(latencies, time.Since(start).Seconds())
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

	fmt.Println("\n=== Benchmark Results ===")
	fmt.Printf("Total successful requests: %d\n", len(latencies))
	fmt.Printf("Average latency: %.6f sec\n", total/float64(len(latencies)))
	fmt.Printf("Max latency: %.6f sec\n", maxLatency)
	fmt.Printf("Min latency: %.6f sec\n", minLatency)
	fmt.Printf("Throughput: %.2f requests/sec\n", float64(len(latencies))/total)
	fmt.Printf("Total benchmark time: %.6f sec\n", time.Since(startBenchmark).Seconds())
}
