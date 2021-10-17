package main

import (
	"fmt"
	. "github.com/craftsangjae/dojocoin/cryto"
)

func main() {
	AddBlock([]*Tx{CoinbaseTx("강상재")})
	AddBlock([]*Tx{CoinbaseTx("강상재")})
	AddBlock([]*Tx{CoinbaseTx("허진수")})
	AddBlock([]*Tx{CoinbaseTx("조서형")})
	AddBlock([]*Tx{CoinbaseTx("선우승환")})

	for _, block := range ListAllBlocks() {
		fmt.Println(block.MarshalJson())
	}
	fmt.Println(GetBlockChain().BalanceByAddress("강상재"))
	fmt.Println(GetBlockChain().BalanceByAddress("허진수"))
}
