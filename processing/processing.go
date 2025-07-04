package processing

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func KillPort(p string) error {
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

func KillMonitorProcesses() {
	out, err := exec.Command("pgrep", "-f", "monitor").Output()
	if err != nil {
		return
	}

	pids := strings.Fields(string(out))
	for _, pid := range pids {
		log.Printf("Killing existing monitor process with PID %s", pid)
		err := exec.Command("kill", "-9", pid).Run()
		if err != nil {
			log.Printf("Failed to kill process %s: %v", pid, err)
		}
	}
}
