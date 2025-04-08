package tcp_test

import (
    "bufio"
    "net"
    "strings"
    "testing"
    "time"
    "solune/tcp"
    "solune/store"
)

func simulateClient(t *testing.T, conn net.Conn, input string, expected string) {
    t.Helper()

    _, err := conn.Write([]byte(input + "\n"))
    if err != nil {
        t.Fatalf("Failed to write to server: %v", err)
    }

    reader := bufio.NewReader(conn)
    response, err := reader.ReadString('\n')
    if err != nil {
        t.Fatalf("Failed to read from server: %v", err)
    }

    if strings.TrimSpace(response) != expected {
        t.Errorf("Expected '%s', got '%s'", expected, strings.TrimSpace(response))
    }
}

func TestHandleClient_Ping(t *testing.T) {
    serverConn, clientConn := net.Pipe()
    defer serverConn.Close()
    defer clientConn.Close()

    manager := store.NewDataStoreManager()
    server := tcp.NewServer(manager)

    go server.HandleClient(serverConn)

    simulateClient(t, clientConn, "PING", "PING")
}

func TestHandleClient_Echo(t *testing.T) {
    serverConn, clientConn := net.Pipe()
    defer serverConn.Close()
    defer clientConn.Close()

    manager := store.NewDataStoreManager()
    server := tcp.NewServer(manager)

    go server.HandleClient(serverConn)

    simulateClient(t, clientConn, "Hello", "Hello")
}

func TestHandleClient_UnknownCommand(t *testing.T) {
    serverConn, clientConn := net.Pipe()
    defer serverConn.Close()
    defer clientConn.Close()

    manager := store.NewDataStoreManager()
    server := tcp.NewServer(manager)

    go server.HandleClient(serverConn)

    simulateClient(t, clientConn, "FOO", "FOO")
}

func TestHandleClient_ClientDisconnect(t *testing.T) {
    serverConn, clientConn := net.Pipe()

    manager := store.NewDataStoreManager()
    server := tcp.NewServer(manager)

    done := make(chan struct{})

    go func() {
        defer close(done)
        server.HandleClient(serverConn)
    }()

    clientConn.Close()

    select {
    case <-done:
    case <-time.After(time.Second):
        t.Fatal("HandleClient did not exit on client disconnect")
    }
}
