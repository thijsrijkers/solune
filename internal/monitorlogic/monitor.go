package monitorlogic

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

const (
	MaxCPUUsage    = 70.0    // percent
	MaxMemoryUsage = 2000.0  // megabytes (2 GB)
)

func Run() {
	for {
		log.Println("Scanning all processes for workers...")

		procs, err := process.Processes()
		if err != nil {
			log.Printf("Failed to list processes: %v", err)
			time.Sleep(10 * time.Minute)
			continue
		}

		var wg sync.WaitGroup

		for _, p := range procs {
			wg.Add(1)
			go func(proc *process.Process) {
				defer wg.Done()
				checkProcess(proc)
			}(p)
		}

		wg.Wait()

		time.Sleep(10 * time.Minute)
	}
}

func checkProcess(p *process.Process) {
	name, err := p.Name()
	if err != nil {
		return
	}

	if !strings.Contains(name, "worker") {
   		return
	}

	pid := p.Pid

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		return
	}

	memMB := float32(memInfo.RSS) / 1024 / 1024

	if cpuPercent > MaxCPUUsage || memMB > MaxMemoryUsage {
		log.Printf("Worker PID %d exceeded limits: CPU %.2f%%, Mem %.2fMB", pid, cpuPercent, memMB)
	}
}

