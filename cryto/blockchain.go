package cryto

import (
	"bytes"
	"crypto/sha256"
	gob "encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/craftsangjae/dojocoin/db"
	"github.com/craftsangjae/dojocoin/utils"
	"strconv"
	"sync"
)

var chain *blockChain
var once sync.Once

func init() {
	GetBlockChain()
}

type block struct {
	Data     string
	Hash     string
	PrevHash string
	Height   int
}

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

func (b *block) MarshalJson() string {
	val, err := json.MarshalIndent(struct {
		Data     string `json:"data"`
		PrevHash string `json:"prevHash"`
		Hash     string `json:"hash"`
	}{
		Data:     b.Data,
		PrevHash: b.PrevHash,
		Hash:     b.Hash,
	}, "", "  ")
	utils.HandleErr(err)
	return string(val)
}

func (chain *blockChain) restore(data []byte) {
	a := gob.NewDecoder(bytes.NewReader(data))
	a.Decode(chain)
}

func (b *block) restore(data []byte) {
	a := gob.NewDecoder(bytes.NewReader(data))
	a.Decode(b)
}

func (b *block) persist() {
	db.SaveBlock(b.Hash, b)
}

func (chain *blockChain) persist() {
	db.SaveBlockchain(chain)
}

// 새 블록을 생성합니다
//
func newBlock(data string) *block {
	prevHash := chain.NewestHash
	newHash := calculateHash(data, prevHash, chain.Height)
	return &block{data, newHash, prevHash, chain.Height + 1}
}

func calculateHash(data string, prevHash string, height int) string {
	hashBytes := sha256.Sum256([]byte(data + prevHash + strconv.Itoa(height)))
	return fmt.Sprintf("%x", hashBytes)
}

func updateChain(b *block) {
	chain.NewestHash = b.Hash
	chain.Height = b.Height
}

func AddBlock(data string) {
	b := newBlock(data)
	b.persist()
	updateChain(b)
	chain.persist()
}

func FindBlock(hash string) *block {
	b := block{"", "", "", 1}
	data := db.FindData(hash)
	b.restore(data)
	return &b
}

func GetBlockChain() *blockChain {
	if chain == nil {
		once.Do(func() {
			checkpoint := db.LoadCheckpoint()
			chain = &blockChain{"", 0}
			if checkpoint == nil {
				AddBlock("GENESIS BLOCK")
			} else {
				chain.restore(checkpoint)
			}
		})
	}
	return chain
}
