package tcp

import (
	"errors"
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
	errInvalidPair   = errors.New("invalid key:value pair")
)

func ParseCommand(command string) (Behavior, error) {
	var b Behavior
	count := 0

	for len(command) > 0 {
		// Find next segment
		seg := command
		if i := strings.IndexByte(command, '|'); i >= 0 {
			seg = command[:i]
			command = command[i+1:]
		} else {
			command = ""
		}

		// Find '=' in segment
		eq := strings.IndexByte(seg, '=')
		if eq < 0 {
			return Behavior{}, errInvalidPair
		}

		key := trimSpace(seg[:eq])
		val := trimSpace(seg[eq+1:])

		switch key {
		case "instruction":
			b.Instruction = val
		case "store":
			b.Store = val
		case "key":
			b.Key = val
		case "data":
			b.Data = val
		default:
			return Behavior{}, errors.New("unknown key: " + key)
		}
		count++
	}

	if count < 2 {
		return Behavior{}, errInvalidFormat
	}

	return b, nil
}

// trimSpace is an inline-friendly replacement for strings.TrimSpace
// that avoids the unicode overhead — TCP commands are ASCII.
func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}
