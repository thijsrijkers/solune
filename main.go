package main

import (
    "solune/store"
    "solune/tcp"
)

func main() {
    manager := store.NewDataStoreManager()
    tcp.StartServer("9000", manager)
}
