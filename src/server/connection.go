package server

import (
    "bufio"
    "fmt"
    "log"
    "net"
)

func HandleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Error reading: %v", err)
            return
        }
        log.Printf("Received: %s", message)

        _, err = fmt.Fprintf(conn, "Echo: %s", message)
        if err != nil {
            log.Printf("Error writing: %v", err)
            return
        }
    }
}
