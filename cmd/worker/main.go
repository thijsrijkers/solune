package main

import (
    "fmt"
    "os"
    "strconv"

    "solune/internal"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <port>")
        os.Exit(1)
    }

    port, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Printf("Invalid port number: %v\n", err)
        os.Exit(1)
    }

    workerlogic.Run(strconv.Itoa(port))
}
