package tcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Behavior struct {
	Instruction string
	Store       string
	Key         string
	Data        string
}

var (
	errInvalidFormat = errors.New("invalid command format")
)

func ParseCommand(command string) (Behavior, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return Behavior{}, errInvalidFormat
	}

	if !strings.HasPrefix(command, "{") {
		return Behavior{}, errors.New("legacy command format removed: send JSON object")
	}

	return parseJSONCommand(command)
}

func parseJSONCommand(command string) (Behavior, error) {
	var b Behavior
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(command), &raw); err != nil {
		return Behavior{}, fmt.Errorf("invalid json command: %w", err)
	}

	for key := range raw {
		switch key {
		case "instruction", "store", "key", "data":
		default:
			return Behavior{}, errors.New("unknown key: " + key)
		}
	}

	rawInstruction, ok := raw["instruction"]
	if !ok {
		return Behavior{}, errInvalidFormat
	}
	if err := json.Unmarshal(rawInstruction, &b.Instruction); err != nil {
		return Behavior{}, fmt.Errorf("instruction must be a string")
	}

	if rawStore, ok := raw["store"]; ok {
		if err := json.Unmarshal(rawStore, &b.Store); err != nil {
			return Behavior{}, fmt.Errorf("store must be a string")
		}
	}

	if rawKey, ok := raw["key"]; ok {
		if err := json.Unmarshal(rawKey, &b.Key); err != nil {
			var number json.Number
			dec := json.NewDecoder(bytes.NewReader(rawKey))
			dec.UseNumber()
			if decErr := dec.Decode(&number); decErr != nil {
				return Behavior{}, fmt.Errorf("key must be a string or number")
			}
			b.Key = number.String()
		}
	}

	if rawData, ok := raw["data"]; ok {
		var str string
		if err := json.Unmarshal(rawData, &str); err != nil {
			return Behavior{}, fmt.Errorf("data must be a string")
		}
		b.Data = str
	}

	return b, nil
}
