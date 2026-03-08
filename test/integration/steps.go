package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const addr = "localhost:9000"

func send(conn net.Conn, command string) string {
	fmt.Fprintf(conn, command+"\n")
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	resp, _ := bufio.NewReader(conn).ReadString('\n')
	return resp
}

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	steps := []struct {
		label   string
		command string
	}{
		{"1. Create store", "instruction=set|store=user_data"},
		{"2. Set data", "instruction=set|store=user_data|data={'name': 'John Doe', 'age': 30}"},
		{"3. Get all data", "instruction=get|store=user_data"},
	}

	for _, step := range steps {
		fmt.Printf("\n[%s]\n> %s\n", step.label, step.command)
		resp := send(conn, step.command)
		fmt.Printf("< %s\n", resp)
		time.Sleep(100 * time.Millisecond)
	}
}
