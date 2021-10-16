package main

import (
	. "github.com/craftsangjae/dojocoin/cryto"
)

func main() {
	GetBlockChain()
	AddBlock("안녕하세요")
	AddBlock("새로운")
	AddBlock("블록체인입니다.")
}
