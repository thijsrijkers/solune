package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func killPort(p string) error {
	// Step 1: Find the PID using lsof
	findCmd := exec.Command("lsof", "-ti", fmt.Sprintf("tcp:%s", p))
	output, err := findCmd.Output()
	if err != nil || len(output) == 0 {
		log.Printf("No process found using TCP port %s", p)
		return nil
	}

	// Step 2: Kill the process
	pid := string(output)
	log.Printf("Killing process using TCP port %s (PID: %s)...", p, pid)
	killCmd := exec.Command("kill", "-9", pid)
	killCmd.Stdout = os.Stdout
	killCmd.Stderr = os.Stderr
	return killCmd.Run()
}

func main() {
	// Step 1: Run the build command
	buildCmd := exec.Command("go", "build", "-o", "worker", "cmd/worker/main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	log.Println("Building worker binary...")
	err := buildCmd.Run()
	if err != nil {
		log.Fatalf("Failed to build worker: %v", err)
	}
	log.Println("Worker binary built successfully.")

	// Step 2: Read db directory
	dbDir := "db"
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		log.Fatalf("Failed to read db directory: %v", err)
	}

	// Step 3: Start a worker for each port, or one if none found
	hasAny := false

	for _, entry := range entries {
		if entry.IsDir() {
			hasAny = true
			port := entry.Name()

			go func(p string) {
				_ = killPort(p) // Kill any process using this port
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
	}

	if !hasAny {
		go func() {
			cmd := exec.Command("./worker", "9000")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			log.Printf("No existing shards found, creating new one....")
			err := cmd.Start()
			if err != nil {
				log.Printf("Failed to start worker with no port: %v", err)
				return
			}

			log.Printf("Worker with newly create shard started with PID %d", cmd.Process.Pid)

			err = cmd.Wait()
			if err != nil {
				log.Printf("Worker with no port exited with error: %v", err)
			} else {
				log.Printf("Worker with no port exited successfully")
			}
		}()
	}

	select {}
}
