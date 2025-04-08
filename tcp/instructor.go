package tcp

import (
	"errors"
	"strings"
)

type Behavior struct {
	Instruction string
	Store       string
}

func ParseCommand(command string) (Behavior, error) {
	parts := strings.Split(command, "|")
	if len(parts) != 2 {
		return Behavior{}, errors.New("invalid command format")
	}

	var b Behavior

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return Behavior{}, errors.New("invalid key:value pair")
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "instruction":
			b.Instruction = value
		case "store":
			b.Store = value
		default:
			return Behavior{}, errors.New("unknown key: " + key)
		}
	}

	return b, nil
}
