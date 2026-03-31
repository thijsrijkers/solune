package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"solune/builder"
	"solune/processing"
	"syscall"
	"time"
)

func main() {
	builder.BuildBinary()
	processing.KillSupervisorProcesses()
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

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigCh)

		go func() {
			sig := <-sigCh
			log.Printf("Forwarding %v to supervisor (PID %d)", sig, supervisorCmd.Process.Pid)
			if err := supervisorCmd.Process.Signal(sig); err != nil {
				log.Printf("Failed to forward signal to supervisor: %v", err)
			}
		}()

		if err := supervisorCmd.Wait(); err != nil {
			log.Printf("Supervisor exited with error: %v", err)
		} else {
			log.Printf("Supervisor exited cleanly")
		}
	}
}
