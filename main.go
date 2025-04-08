package main

import (
    "solune/store"
    "solune/tcp"
)

func main() {
    manager := store.NewDataStoreManager()
    manager.AddStore("users")

	usersStore, _ := manager.GetStore("users")
	usersStore.Set("1", map[string]interface{}{"name": "root", "password": ""})

    tcp.StartServer("9000", manager)
}
