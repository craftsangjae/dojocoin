package cryto

import (
	"github.com/craftsangjae/dojocoin/utils"
	"time"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string
	Timestamp int
	TxIns     []*TxIn
	TxOuts    []*TxOut
}

func NewTx(txIns []*TxIn, txOuts []*TxOut) *Tx {
	tx := Tx{"", int(time.Now().Unix()), txIns, txOuts}
	tx.Id = utils.Hash(tx)
	return &tx
}

type TxIn struct {
	TxID  string
	Index int
	Owner string
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

func CoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	return NewTx(txIns, txOuts)
}
