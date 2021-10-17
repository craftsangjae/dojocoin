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
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

func CoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	return NewTx(txIns, txOuts)
}
