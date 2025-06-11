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
	wg.Add(2)

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

	wg.Wait()
	log.Println("All binaries built successfully.")
}

