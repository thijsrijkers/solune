package builder

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func BuildBinary() {
	log.Println("Building binaries...")

	var wg sync.WaitGroup
	wg.Add(4)

	// Build worker binary
	go func() {
		defer wg.Done()
		log.Println("Building worker binary...")
		cmd := exec.Command("go", "build", "-o", "worker", "cmd/worker/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build worker: %v", err)
		}
		log.Println("Worker binary built successfully.")
	}()

	// Build supervisor binary
	go func() {
		defer wg.Done()
		log.Println("Building supervisor binary...")
		cmd := exec.Command("go", "build", "-o", "supervisor", "cmd/supervisor/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build supervisor: %v", err)
		}
		log.Println("Supervisor binary built successfully.")
	}()

	// Build monitor binary
	go func() {
		defer wg.Done()
		log.Println("Building monitor binary...")
		cmd := exec.Command("go", "build", "-o", "monitor", "cmd/monitor/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build rebalancing: %v", err)
		}
		log.Println("Monitor binary built successfully.")
	}()

	// Build rebalancing binary
	go func() {
		defer wg.Done()
		log.Println("Building rebalancing binary...")
		cmd := exec.Command("go", "build", "-o", "rebalancing", "cmd/rebalancing/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build rebalancing: %v", err)
		}
		log.Println("Rebalancing binary built successfully.")
	}()

	wg.Wait()
	log.Println("All binaries built successfully.")
}

