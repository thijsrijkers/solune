package builder

import (
	"log"
	"os"
	"os/exec"
)

func BuildBinary() {
	log.Println("Building worker binary...")
	
	buildCmd := exec.Command("go", "build", "-o", "worker", "cmd/worker/main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build worker: %v", err)
	}
	
	log.Println("Worker binary built successfully.")

	log.Println("Building supervisor binary...")
	
	buildCmd = exec.Command("go", "build", "-o", "supervisor", "cmd/supervisor/main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build supervisor: %v", err)
	}
	
	log.Println("Supervisor binary built successfully.")
}
