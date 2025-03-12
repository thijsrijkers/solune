package main

import (
    "stack/src/store"
    "stack/src/tcp"
)

func main() {
    manager := store.NewDataStoreManager()
    tcp.StartServer("9000", manager)
}
