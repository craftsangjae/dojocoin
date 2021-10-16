package cryto

import (
	"crypto/sha256"
	"fmt"
	"github.com/craftsangjae/dojocoin/db"
	"strconv"
	"sync"
)

var chain *blockChain
var once sync.Once

func init() {
	GetBlockChain()
}

type block struct {
	data     string
	hash     string
	prevHash string
	height   int `json:"height"`
}

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

func (b *block) persist() {
	db.SaveBlock(b.hash, b)
}

func (b *blockChain) persist() {
	db.SaveBlockchain(b)
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
	chain.NewestHash = b.hash
	chain.Height = b.height
}

func AddBlock(data string) {
	b := newBlock(data)
	b.persist()
	updateChain(b)
	chain.persist()
}

func GetBlockChain() *blockChain {
	if chain == nil {
		once.Do(func() {
			chain = &blockChain{}
			AddBlock("GENESIS BLOCK")
		})
	}
	return chain
}
