package data

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register("")       // string
	gob.Register(0)        // int
	gob.Register(true)     // bool
	gob.Register(float32(0)) // float32
	gob.Register(float64(0)) // float64
}

func MapToBinary(m map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func BinaryToMap(data []byte) (map[string]interface{}, error) {
	buf := bytes.NewReader(data)

	dec := gob.NewDecoder(buf)
	var m map[string]interface{}

	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
