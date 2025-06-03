package main

import (
	"log"
	"os"
	"os/exec"
)

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

	// Step 2: Start the worker binary
	cmd := exec.Command("./worker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Starting worker process...")
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start worker process: %v", err)
	}

	log.Println("Worker process started with PID:", cmd.Process.Pid)

	// Optional: wait for worker to exit
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Worker process exited with error: %v", err)
	}
}
