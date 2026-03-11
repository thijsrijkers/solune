package main

import (
	"log"
	"os"
	"os/exec"
	"solune/builder"
	"solune/processing"
	"time"
)

func main() {
	builder.BuildBinary()
	processing.KillMonitorProcesses()
	port := "9000"
	processing.KillPort(port)
	time.Sleep(500 * time.Millisecond)

	supervisorCmd := exec.Command("./supervisor", port)
	supervisorCmd.Stdout = os.Stdout
	supervisorCmd.Stderr = os.Stderr
	err := supervisorCmd.Start()

	if err != nil {
		log.Printf("Failed to start supervisor for port %s: %v", port, err)
	} else {
		log.Printf("Supervisor for port %s started with PID %d", port, supervisorCmd.Process.Pid)
	}
}
