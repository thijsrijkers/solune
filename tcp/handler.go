package tcp

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "solune/sql"
    "solune/store"
)

type Server struct {
    sqlTranslator *sql.SQLDataTranslator
}

func NewServer(manager *store.DataStoreManager) *Server {
    return &Server{sqlTranslator: sql.NewSQLDataTranslator(manager)}
}

func (s *Server) HandleClient(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

    msg = strings.TrimSpace(msg)
	fmt.Println("Processed Query:", msg)

    result, err := s.sqlTranslator.Translate(msg)
	if err != nil {
		fmt.Println("Error translating SQL:", err)
		conn.Write([]byte("ERROR: " + err.Error() + "\n"))
		return
	}

	conn.Write([]byte(fmt.Sprintf("Query result: %v\n", result)))
}
