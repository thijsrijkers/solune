package tcp

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

func HandleReadError(err error) {
	if err != io.EOF {
		log.Println("Read error:", err)
	}
}

func WriteError(writer *bufio.Writer, err error) {
	resp := map[string]interface{}{"error": err.Error()}
	jsonData, _ := json.Marshal(resp)
	writer.WriteString(string(jsonData) + "\n")
	writer.Flush()
}

func WriteResult(writer *bufio.Writer, result []map[string]interface{}) {
	if len(result) == 0 {
		writer.WriteString("{\"status\":404}\n")
		writer.Flush()
		return
	}

	if data, ok := result[0]["data"].(string); ok {
		writer.WriteString(data)
		writer.WriteByte('\n')
		writer.Flush()
		return
	}

	writer.WriteString("{\"status\":200}\n")
	writer.Flush()
}
