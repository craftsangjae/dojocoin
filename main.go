package main

import (
	"fmt"
	. "github.com/craftsangjae/dojocoin/cryto"
)

func main() {
	AddBlock("안녕하세요")
	AddBlock("새로운")
	AddBlock("형태의")
	AddBlock("블록체인입니다.")

	for _, block := range ListAllBlocks() {
		fmt.Println(block.MarshalJson())
	}
}
