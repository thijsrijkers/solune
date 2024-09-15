package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
)

// handleConnection manages the communication for each TCP connection
func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Read data from the connection
    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Error reading: %v", err)
            return
        }
        log.Printf("Received: %s", message)

        // Echo the message back to the client
        _, err = conn.Write([]byte(fmt.Sprintf("Echo: %s", message)))
        if err != nil {
            log.Printf("Error writing: %v", err)
            return
        }
    }
}
