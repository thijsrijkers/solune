package main

import (
    "log"
    "solune/internal/rebalancinglogic"
)

func main() {

    log.Printf("Starting rebalancing...")
    rebalancinglogic.Run()
}

