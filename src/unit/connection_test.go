package unit

import (
    "bytes"
    "testing"
    "paper/src/unit/fixtures"
)

func TestHandleConnectionMock(t *testing.T) {
    input := bytes.NewBufferString("Hello, Server!\n")
    output := &bytes.Buffer{}

    fixtures.HandleConnectionMock(input, output)

    expected := "Echo: Hello, Server!\n"
    if output.String() != expected {
        t.Errorf("Expected '%s' but got '%s'", expected, output.String())
    }
}
