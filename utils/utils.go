package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func ToBytes(i interface{}) []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	err := encoder.Encode(i)
	if err != nil {
		return nil
	}
	return blockBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	a := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(a.Decode(i))
}
