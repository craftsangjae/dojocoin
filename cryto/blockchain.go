package cryto

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
	height int `json:"height"`
}

type blockChain struct {
	blocks []*block
}


// 새 블록을 생성합니다
//
func newBlock(data string, prevHash string) *block {
	newHash := calculateHash(data, prevHash)
	return &block{data ,newHash,prevHash, len(chain.blocks)}
}

func calculateHash(data string, prevHash string) string {
	hashBytes := sha256.Sum256([]byte(data + prevHash))
	return fmt.Sprintf("%x", hashBytes)
}

func isEmpty() bool {
	if chain == nil {
		return true
	}
	if len(chain.blocks) == 0 {
		return true
	}
	return false
}

func (chain *blockChain) getLastHash() string {
	if isEmpty() {
		return ""
	}
	return chain.blocks[len(chain.blocks) - 1].hash
}

func AddBlock(data string) {
	chain.blocks = append(chain.blocks, newBlock(data, chain.getLastHash()))
}

func ListBlocks() []*block {
	return chain.blocks
}

func FindBlock(height int) *block {
	return chain.blocks[height]
}

func (b *block) MarshalJson() ([]byte, error) {
	j, err := json.MarshalIndent(struct {
		Data string `json:"data"`
		Hash string `json:"hash"`
		PrevHash string `json:"prevHash,omitempty"`
	}{
		Data: b.data,
		Hash: b.hash,
		PrevHash: b.prevHash,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	return j, nil
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