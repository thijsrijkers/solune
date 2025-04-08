package tcp

import (
    "fmt"
    "net"
    "solune/store"
)

func StartServer(port string, manager *store.DataStoreManager) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Listening on port", port)

    server := NewServer(manager)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }
        go server.HandleClient(conn)
    }
}
