package main

import (
	"fmt"
	. "github.com/craftsangjae/dojocoin/cryto"
)

func main() {
	AddBlock("안녕하세요")
	AddBlock("새로운")
	AddBlock("블록체인입니다.")

	for _, block := range ListBlocks() {
		if value, err := block.MarshalJson(); err == nil {
			fmt.Println(string(value))
		}
	}
}