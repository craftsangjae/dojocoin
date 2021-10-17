package blockchain

import (
	"errors"
)

var Mempool = &mempool{}

type mempool struct {
	Txs []*Tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from) < amount {
		return nil, errors.New("NOT ENOUGH MONEY")
	}

	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	for _, utxOut := range UTxOutsByAddress(from) {
		if total > amount {
			break
		}
		txIns = append(txIns, &TxIn{utxOut.TxID, utxOut.Index, from})
		total += utxOut.Amount
	}

	if total > amount {
		txOuts = append(txOuts, &TxOut{from, total - amount})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	return NewTx(txIns, txOuts), nil
}

func isOnMempool(utxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == utxOut.TxID && input.Index == utxOut.Index {
				return true
			}
		}
	}
	return false
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("강상재", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := CoinbaseTx("강상재")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
