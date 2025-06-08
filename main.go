package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"solune/tcprelay"
)

func killPort(p string) error {
	findCmd := exec.Command("lsof", "-ti", fmt.Sprintf("tcp:%s", p))
	output, err := findCmd.Output()
	if err != nil || len(output) == 0 {
		log.Printf("No process found using TCP port %s", p)
		return nil
	}

	pid := strings.TrimSpace(string(output))
	log.Printf("Killing process using TCP port %s (PID: %s)...", p, pid)
	killCmd := exec.Command("kill", "-9", pid)
	killCmd.Stdout = os.Stdout
	killCmd.Stderr = os.Stderr
	return killCmd.Run()
}

func main() {
	log.Println("Building worker binary...")
	buildCmd := exec.Command("go", "build", "-o", "worker", "cmd/worker/main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build worker: %v", err)
	}
	log.Println("Worker binary built successfully.")

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

	// Step 1: Launch all workers
	for _, port := range allPorts {
		go func(p string) {
			_ = killPort(p)
			time.Sleep(500 * time.Millisecond)

			cmd := exec.Command("./worker", p)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			log.Printf("Starting worker for port %s...", p)
			err := cmd.Start()
			if err != nil {
				log.Printf("Failed to start worker for port %s: %v", p, err)
				return
			}

			log.Printf("Worker for port %s started with PID %d", p, cmd.Process.Pid)
			err = cmd.Wait()
			if err != nil {
				log.Printf("Worker for port %s exited with error: %v", p, err)
			} else {
				log.Printf("Worker for port %s exited successfully", p)
			}
		}(port)
	}

	// Step 2: Launch this node's TCP relay on 8743
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

