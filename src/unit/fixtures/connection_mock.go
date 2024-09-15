package fixtures

import (
    "bufio"
    "fmt"
    "io"
)

func HandleConnectionMock(reader io.Reader, writer io.Writer) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        message := scanner.Text()
        _, _ = fmt.Fprintf(writer, "Echo: %s\n", message)
    }
}
