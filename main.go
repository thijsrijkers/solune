package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"solune/tcprelay"
	"solune/processing"
	"solune/builder"
)

func main() {
	//Step 1: Setup binary, cleanup old monitor process and retrieve argument data
	builder.BuildBinary();

	processing.KillMonitorProcesses();

	dbDir := "db"
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		log.Fatalf("Failed to read db directory: %v", err)
	}

	var allPorts []string
	for _, entry := range entries {
		if entry.IsDir() {
			allPorts = append(allPorts, entry.Name())
		}
	}
	if len(allPorts) == 0 {
		allPorts = append(allPorts, "9000")
	}

	// Step 2: Launch all workers + supervisors
	for _, port := range allPorts {
		go func(p string) {
			_ = processing.KillPort(p)
			time.Sleep(500 * time.Millisecond)

			// Start worker
			cmd := exec.Command("./worker", p)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			log.Printf("Starting worker for port %s...", p)
			err := cmd.Start()
			if err != nil {
				log.Printf("Failed to start worker for port %s: %v", p, err)
				return
			}

			workerPid := cmd.Process.Pid
			log.Printf("Worker for port %s started with PID %d", p, workerPid)

			// Start supervisor
			supervisorCmd := exec.Command("./supervisor", p, fmt.Sprintf("%d", workerPid))
			supervisorCmd.Stdout = os.Stdout
			supervisorCmd.Stderr = os.Stderr

			log.Printf("Starting supervisor for PID %d on port %s", workerPid, p)

			err = supervisorCmd.Start()
			if err != nil {
				log.Printf("Failed to start supervisor for port %s: %v", p, err)
			} else {
				log.Printf("Supervisor for port %s started with PID %d", p, supervisorCmd.Process.Pid)
			}

			err = cmd.Wait()
			if err != nil {
				log.Printf("Worker for port %s exited with error: %v", p, err)
			} else {
				log.Printf("Worker for port %s exited successfully", p)
			}
		}(port)
	}

	// Step 3: Launch monitor process
	monitorCmd := exec.Command("./monitor")
	monitorCmd.Stdout = os.Stdout
	monitorCmd.Stderr = os.Stderr
	monitorStartErr := monitorCmd.Start()
	if monitorStartErr != nil {
		log.Fatalf("Failed to start monitor: %v", monitorStartErr)
	}
	log.Printf("Started monitor with PID %d", monitorCmd.Process.Pid)

	// Step 4: Launch this node's TCP relay on 8743
	relayPort := "8743"

	var peers []string
	for _, p := range allPorts {
		if p != relayPort {
			peers = append(peers, p)
		}
	}

	node := tcprelay.NewRelayNode(relayPort, peers)
	err = node.Start()
	if err != nil {
		log.Fatalf("Failed to start TCP relay: %v", err)
	}
}

