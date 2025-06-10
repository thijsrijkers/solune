package main

import (
    "fmt"
    "os"
    "solune/internal/supervisorlogic"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run main.go <port> <pid>")
        os.Exit(1)
    }

    port := os.Args[1]
    pid := os.Args[2]

    supervisorlogic.Run(port, pid)
}
