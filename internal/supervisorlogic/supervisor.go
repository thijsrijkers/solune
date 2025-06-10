package supervisorlogic

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"solune/processing"
)

func Run(port string, pid string) {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Fatalf("Invalid PID provided: %v", err)
	}

	log.Printf("Starting supervisor for PID %d on port %s", pidInt, port)

	for {
		if !isProcessRunning(pidInt) {
			log.Printf("Process %d stopped. Restarting shard on port %s...", pidInt, port)

			err := processing.KillPort(port)
			if err != nil {
				log.Printf("Error killing port %s: %v", port, err)
			}

			proc, err := startShard(port)
			if err != nil {
				log.Printf("Failed to start new shard on port %s: %v", port, err)
				time.Sleep(5 * time.Second)
				continue
			}

			pidInt = proc.Pid
			log.Printf("Now supervising new shard with PID %d on port %s", pidInt, port)
		}

		time.Sleep(60 * time.Second)
	}
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
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

func startShard(port string) (*os.Process, error) {
	_ = killPort(port)
	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command("./worker", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting worker for port %s...", port)
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start worker for port %s: %v", port, err)
		return nil, err
	}

	log.Printf("Worker for port %s started with PID %d", port, cmd.Process.Pid)

	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Printf("Worker for port %s exited with error: %v", port, err)
		} else {
			log.Printf("Worker for port %s exited successfully", port)
		}
	}()

	return cmd.Process, nil
}
