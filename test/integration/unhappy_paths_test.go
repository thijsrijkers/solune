package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func sendCommandResult(conn net.Conn, command string) (map[string]interface{}, error) {
	if _, err := fmt.Fprintf(conn, command+"\n"); err != nil {
		return nil, fmt.Errorf("failed to send command %q: %w", command, err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response for command %q: %w", command, err)
	}

	line = strings.TrimSpace(line)
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		return nil, fmt.Errorf("response is not valid JSON for command %q: %q (%w)", command, line, err)
	}

	return payload, nil
}

func sendCommand(t *testing.T, conn net.Conn, command string) map[string]interface{} {
	t.Helper()

	payload, err := sendCommandResult(conn, command)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return payload
}

func expectStatus(t *testing.T, payload map[string]interface{}, expected int) {
	t.Helper()

	status, ok := payload["status"].(float64)
	if !ok {
		t.Fatalf("expected status response, got: %#v", payload)
	}

	if int(status) != expected {
		t.Fatalf("expected status %d, got %d in response %#v", expected, int(status), payload)
	}
}

func expectErrorContains(t *testing.T, conn net.Conn, command, expectedFragment string) {
	t.Helper()

	payload := sendCommand(t, conn, command)
	errMsg, ok := payload["error"].(string)
	if !ok {
		t.Fatalf("expected error response for command %q, got: %#v", command, payload)
	}

	if !strings.Contains(errMsg, expectedFragment) {
		t.Fatalf("expected error to contain %q for command %q, got %q", expectedFragment, command, errMsg)
	}
}

func TestIntegrationUnhappyPaths(t *testing.T) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("failed to connect to %s: %v", addr, err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	})

	storeName := fmt.Sprintf("unhappy_paths_%d", time.Now().UnixNano())
	missingStore := fmt.Sprintf("missing_%d", time.Now().UnixNano())

	setup := sendCommand(t, conn, fmt.Sprintf("instruction=set|store=%s", storeName))
	expectStatus(t, setup, 200)

	t.Cleanup(func() {
		cleanupConn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Logf("cleanup failed to connect for store %s: %v", storeName, err)
			return
		}
		defer cleanupConn.Close()

		payload, err := sendCommandResult(cleanupConn, fmt.Sprintf("instruction=delete|store=%s", storeName))
		if err != nil {
			t.Logf("cleanup failed for store %s: %v", storeName, err)
			return
		}

		if status, ok := payload["status"].(float64); ok && int(status) == 200 {
			return
		}

		if errMsg, ok := payload["error"].(string); ok && strings.Contains(errMsg, "not found") {
			return
		}

		t.Logf("cleanup returned unexpected payload for store %s: %#v", storeName, payload)
	})

	testCases := []struct {
		name             string
		command          string
		expectedFragment string
	}{
		{
			name:             "invalid pair when '=' missing",
			command:          "instruction",
			expectedFragment: "invalid key:value pair",
		},
		{
			name:             "invalid pair in second segment",
			command:          "instruction=get|store",
			expectedFragment: "invalid key:value pair",
		},
		{
			name:             "unknown key",
			command:          "unknown=value",
			expectedFragment: "unknown key: unknown",
		},
		{
			name:             "unsupported instruction",
			command:          "instruction=patch",
			expectedFragment: "unsupported action: patch",
		},
		{
			name:             "set with no store",
			command:          "instruction=set",
			expectedFragment: "failed to set data: no store provided",
		},
		{
			name:             "set with existing store but no key and no data",
			command:          fmt.Sprintf("instruction=set|store=%s", storeName),
			expectedFragment: "failed to set data: panic return",
		},
		{
			name:             "set with invalid key",
			command:          fmt.Sprintf("instruction=set|store=%s|key=abc|data=value", storeName),
			expectedFragment: "invalid integer key 'abc'",
		},
		{
			name:             "get from missing store",
			command:          fmt.Sprintf("instruction=get|store=%s", missingStore),
			expectedFragment: fmt.Sprintf("store '%s' not found", missingStore),
		},
		{
			name:             "get with invalid key",
			command:          fmt.Sprintf("instruction=get|store=%s|key=abc", storeName),
			expectedFragment: "invalid integer key 'abc'",
		},
		{
			name:             "get missing key",
			command:          fmt.Sprintf("instruction=get|store=%s|key=999", storeName),
			expectedFragment: "key 999 not found",
		},
		{
			name:             "delete with no store",
			command:          "instruction=delete",
			expectedFragment: "failed to remove store: no store provided",
		},
		{
			name:             "delete missing store",
			command:          fmt.Sprintf("instruction=delete|store=%s", missingStore),
			expectedFragment: fmt.Sprintf("store '%s' not found", missingStore),
		},
		{
			name:             "delete with invalid key",
			command:          fmt.Sprintf("instruction=delete|store=%s|key=abc", storeName),
			expectedFragment: "invalid integer key 'abc'",
		},
		{
			name:             "delete missing key",
			command:          fmt.Sprintf("instruction=delete|store=%s|key=999", storeName),
			expectedFragment: "failed to delete key 999: key 999 not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectErrorContains(t, conn, tc.command, tc.expectedFragment)
		})
	}
}
