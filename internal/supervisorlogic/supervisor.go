package supervisorlogic

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"solune/processing"
	"syscall"
	"strings"
	"time"
)

func Run(port string) {
	log.Printf("Supervisor starting worker on port %s...", port)

	done, err := startWorker(port)
	if err != nil {
		log.Fatalf("Failed to start initial worker on port %s: %v", port, err)
	}

	for {
		exitErr := <-done
		log.Printf("Worker on port %s stopped: %v. Restarting...", port, exitErr)

		if err := processing.KillPort(port); err != nil {
			log.Printf("Error killing port %s: %v", port, err)
		}

		done, err = startWorker(port)
		if err != nil {
			log.Printf("Failed to restart worker on port %s: %v", port, err)
			time.Sleep(5 * time.Second)
			done = make(chan error, 1)
			done <- fmt.Errorf("start failed, retrying")
		}
	}
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}

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

func startWorker(port string) (chan error, error) {
	_ = killPort(port)
	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command("./worker", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting worker for port %s...", port)
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start worker for port %s: %v", port, err)
		return nil, err
	}

	log.Printf("Worker for port %s started with PID %d", port, cmd.Process.Pid)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	return done, nil
}
