package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"solune/builder"
	"solune/processing"
	"time"
)

func main() {
	//Step 1: Setup binary, cleanup old monitor process and retrieve argument data
	builder.BuildBinary()

	processing.KillMonitorProcesses()

	port := "9000"

	processing.KillPort(port)
	time.Sleep(500 * time.Millisecond)

	// Start worker
	cmd := exec.Command("./worker", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting worker for port %s...", port)
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start worker for port %s: %v", port, err)
		return
	}

	workerPid := cmd.Process.Pid
	log.Printf("Worker for port %s started with PID %d", port, workerPid)

	// Start supervisor
	supervisorCmd := exec.Command("./supervisor", port, fmt.Sprintf("%d", workerPid))
	supervisorCmd.Stdout = os.Stdout
	supervisorCmd.Stderr = os.Stderr

	log.Printf("Starting supervisor for PID %d on port %s", workerPid, port)

	err = supervisorCmd.Start()
	if err != nil {
		log.Printf("Failed to start supervisor for port %s: %v", port, err)
	} else {
		log.Printf("Supervisor for port %s started with PID %d", port, supervisorCmd.Process.Pid)
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Worker for port %s exited with error: %v", port, err)
	} else {
		log.Printf("Worker for port %s exited successfully", port)
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
}
