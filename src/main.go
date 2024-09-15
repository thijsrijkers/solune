package main

import (
    "log"
    "net"
    "paper/src/server"
)

func main() {
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Error starting TCP server: %v", err)
    }
    defer ln.Close()

    log.Println("Server is listening on port 8080...")

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v", err)
            continue
        }
        log.Printf("Connection established with %s", conn.RemoteAddr())

        go server.HandleConnection(conn)
    }
}