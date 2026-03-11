package main

import (
    "fmt"
    "os"
    "solune/internal/supervisorlogic"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <port>")
        os.Exit(1)
    }

    port := os.Args[1]

    supervisorlogic.Run(port)
}
