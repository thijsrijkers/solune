package supervisorlogic

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"solune/processing"
	"strings"
	"syscall"
	"time"
)

func Run(port string) {
	log.Printf("Supervisor starting worker on port %s...", port)

	cmd, done, err := startWorker(port)
	if err != nil {
		log.Fatalf("Failed to start initial worker on port %s: %v", port, err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	for {
		select {
		case sig := <-sigCh:
			log.Printf("Supervisor received %v, stopping worker on port %s", sig, port)
			if cmd != nil && cmd.Process != nil {
				_ = cmd.Process.Signal(syscall.SIGTERM)
				select {
				case <-done:
				case <-time.After(2 * time.Second):
					_ = cmd.Process.Kill()
				}
			}
			if err := processing.KillPort(port); err != nil {
				log.Printf("Error cleaning up port %s during shutdown: %v", port, err)
			}
			return

		case exitErr := <-done:
			log.Printf("Worker on port %s stopped: %v. Restarting...", port, exitErr)

			if err := processing.KillPort(port); err != nil {
				log.Printf("Error killing port %s: %v", port, err)
			}

			cmd, done, err = startWorker(port)
			if err != nil {
				log.Printf("Failed to restart worker on port %s: %v", port, err)
				time.Sleep(5 * time.Second)
				done = make(chan error, 1)
				done <- fmt.Errorf("start failed, retrying")
			}
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

func startWorker(port string) (*exec.Cmd, chan error, error) {
	_ = killPort(port)
	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command("./worker", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting worker for port %s...", port)
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start worker for port %s: %v", port, err)
		return nil, nil, err
	}

	log.Printf("Worker for port %s started with PID %d", port, cmd.Process.Pid)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	return cmd, done, nil
}
