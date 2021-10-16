package main

import (
	"fmt"
	. "github.com/craftsangjae/dojocoin/cryto"
)

func main() {
	chain := GetBlockChain()
	AddBlock("안녕하세요")
	AddBlock("새로운")
	AddBlock("블록체인입니다.")

	fmt.Println(chain)
	currentHash := chain.NewestHash
	for true {
		if currentHash == "" {
			break
		}
		block := FindBlock(currentHash)
		currentHash = block.PrevHash
		fmt.Println(block.MarshalJson())
	}
}
