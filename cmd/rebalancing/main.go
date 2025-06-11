package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"

    "solune/internal/rebalancinglogic"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: rebalancing <pid1> <pid2> ...")
        os.Exit(1)
    }

    var pids []int
    for _, arg := range os.Args[1:] {
        pid, err := strconv.Atoi(arg)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid PID: %s\n", arg)
            os.Exit(1)
        }
        pids = append(pids, pid)
    }

    fmt.Printf("Starting rebalancing for PIDs: %s\n", strings.Join(os.Args[1:], ", "))
    rebalancinglogic.Run(pids)
}

