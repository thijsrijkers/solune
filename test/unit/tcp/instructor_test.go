package tcp_test

import (
	"testing"

	"solune/tcp"
)

func TestParseJSONCommand(t *testing.T) {
	t.Run("parses JSON with stringified data", func(t *testing.T) {
		cmd := `{"instruction":"set","store":"users","key":1,"data":"{\"name\":\"John\",\"age\":30}"}`
		parsed, err := tcp.ParseCommand(cmd)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if parsed.Instruction != "set" {
			t.Fatalf("expected instruction 'set', got %q", parsed.Instruction)
		}
		if parsed.Store != "users" {
			t.Fatalf("expected store 'users', got %q", parsed.Store)
		}
		if parsed.Key != "1" {
			t.Fatalf("expected key '1', got %q", parsed.Key)
		}
		if parsed.Data != `{"name":"John","age":30}` {
			t.Fatalf("unexpected data payload: %q", parsed.Data)
		}
	})

	t.Run("parses JSON with string data", func(t *testing.T) {
		cmd := `{"instruction":"set","store":"users","data":"hello"}`
		parsed, err := tcp.ParseCommand(cmd)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if parsed.Data != "hello" {
			t.Fatalf("expected data 'hello', got %q", parsed.Data)
		}
	})

	t.Run("rejects unknown keys", func(t *testing.T) {
		cmd := `{"instruction":"get","unknown":"x"}`
		_, err := tcp.ParseCommand(cmd)
		if err == nil {
			t.Fatalf("expected error for unknown key")
		}
	})

	t.Run("rejects missing instruction", func(t *testing.T) {
		cmd := `{"store":"users"}`
		_, err := tcp.ParseCommand(cmd)
		if err == nil {
			t.Fatalf("expected error for missing instruction")
		}
	})

	t.Run("rejects non-string data", func(t *testing.T) {
		cmd := `{"instruction":"set","store":"users","data":{"name":"John"}}`
		_, err := tcp.ParseCommand(cmd)
		if err == nil {
			t.Fatalf("expected error for non-string data")
		}
	})

	t.Run("rejects legacy format", func(t *testing.T) {
		cmd := "instruction=set|store=users"
		_, err := tcp.ParseCommand(cmd)
		if err == nil {
			t.Fatalf("expected error for legacy format")
		}
	})
}
