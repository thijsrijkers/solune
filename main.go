package main

import (
    "paper/src/store"
    "paper/src/tcp"
)

func main() {
    manager := store.NewDataStoreManager()
    tcp.StartServer("9000", manager)
}
