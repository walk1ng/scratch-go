package rpc

import (
	"bytes"
	"encoding/gob"
)

type RPCData struct {
	Name string
	Args []interface{}
}

func encode(data RPCData) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(b []byte) (*RPCData, error) {
	buf := bytes.NewBuffer(b)
	var data RPCData
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
