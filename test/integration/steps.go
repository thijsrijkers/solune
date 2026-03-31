package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const addr = "127.0.0.1:9000"

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
		os.Exit(1)
	}
	defer conn.Close()

	steps := []struct {
		label   string
		command string
	}{
		{"1. Create store", "instruction=set|store=user_data"},
		{"2. Set data", "instruction=set|store=user_data|data={'name': 'John Doe', 'age': 30}"},
		{"3. Get all stores", "instruction=get"},
		{"4. Get all data", "instruction=get|store=user_data"},
		{"5. Get data before adjustment", "instruction=get|store=user_data|key=1"},
		{"6. Adjust data", "instruction=set|store=user_data|key=1|data={'name': 'John Not Doe', 'age': 0}"},
		{"7. Get data after adjustment", "instruction=get|store=user_data|key=1"},
		{"8. Delete data", "instruction=delete|store=user_data|key=1"},
		{"9. Delete store", "instruction=delete|store=user_data"},
	}

	for _, step := range steps {
		fmt.Printf("\n[%s]\n> %s\n", step.label, step.command)
		resp := send(conn, step.command)
		fmt.Printf("< %s\n", resp)
		time.Sleep(100 * time.Millisecond)
	}
}
