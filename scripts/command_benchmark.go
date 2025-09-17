package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	serverIP    = "127.0.0.1"
	serverPort  = "8743"
	numRequests = 100
	command     = ""
)

func main() {
	latencies := make([]float64, numRequests)
	
	startBenchmark := time.Now()
	
	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	reader := bufio.NewReader(conn)
	
    for i := 0; i < numRequests; i++ {
        start := time.Now()

        _, err := fmt.Fprintf(conn, "%s\n", command)
        if err != nil {
            fmt.Println("Send error:", err)
            continue
        }

        _, err = reader.ReadString('\n')
        if err != nil {
            fmt.Println("Read error:", err)
            continue
        }

        latencies[i] = time.Since(start).Seconds()
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
	fmt.Printf("Total requests: %d\n", numRequests)
	fmt.Printf("Average latency: %.6f sec\n", total/float64(len(latencies)))
	fmt.Printf("Max latency: %.6f sec\n", maxLatency)
	fmt.Printf("Min latency: %.6f sec\n", minLatency)
	fmt.Printf("Throughput: %.2f requests/sec\n", float64(numRequests)/total)
	fmt.Printf("Total benchmark time: %.6f sec\n", time.Since(startBenchmark).Seconds())
}
