package main

import (
	"fmt"
	. "github.com/craftsangjae/dojocoin/blockchain"
	"os"
)

func main() {
	Mempool.AddTx("진수", 10)
	Mempool.AddTx("서형", 10)
	Mempool.AddTx("진수", 10)
	AddBlock()
	for _, name := range []string{"강상재", "진수", "서형"} {
		fmt.Printf("%s : %d\n", name, BalanceByAddress(name))
	}
	fmt.Println("------------------------")

	Mempool.AddTx("서형", 10)
	AddBlock()
	for _, name := range []string{"강상재", "진수", "서형"} {
		fmt.Printf("%s : %d\n", name, BalanceByAddress(name))
	}
	fmt.Println("------------------------")

	os.Remove("./blockchain.db")
}
