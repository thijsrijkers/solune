package tcp

import (
	"bufio"
	"bytes"
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
		writer.WriteString(`{"status":404}` + "\n")
		writer.Flush()
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, item := range result {
		if err := enc.Encode(item); err != nil {
			log.Println("Error encoding item:", err)
			buf.WriteString(`{"error":"failed to serialize"}` + "\n")
		}
	}

	writer.Write(buf.Bytes())
	writer.Flush()
}
