package cryto

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/craftsangjae/dojocoin/db"
	"github.com/craftsangjae/dojocoin/utils"
	"strconv"
	strings "strings"
	"sync"
)

const (
	difficulty int = 2
)

var chain *blockChain
var once sync.Once

func init() {
	GetBlockChain()
}

type block struct {
	Data       string
	Hash       string
	PrevHash   string
	Height     int
	Difficulty int
	Nonce      int
}

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

func (b *block) MarshalJson() string {
	val, err := json.MarshalIndent(b, "", "  ")
	utils.HandleErr(err)
	return string(val)
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

	target := createTarget()
	nounce := 1
	for {
		newHash := calculateHash(data, prevHash, chain.Height, nounce)
		if strings.HasPrefix(newHash, target) {
			return &block{data, newHash, prevHash, chain.Height + 1, difficulty, nounce}
		}
		nounce++
	}
}

func createTarget() string {
	return strings.Repeat("0", difficulty)
}

func calculateHash(data string, prevHash string, height int, nounce int) string {
	hashBytes := sha256.Sum256([]byte(data + prevHash + strconv.Itoa(height) + strconv.Itoa(nounce)))
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
	b := &block{"", "", "", 1, difficulty, 0}
	data := db.FindData(hash)
	utils.FromBytes(b, data)
	return b
}

func GetBlockChain() *blockChain {
	if chain == nil {
		once.Do(func() {
			checkpoint := db.LoadCheckpoint()
			chain = &blockChain{"", 0}
			if checkpoint == nil {
				AddBlock("GENESIS BLOCK")
			} else {
				utils.FromBytes(chain, checkpoint)
			}
		})
	}
	return chain
}

func ListAllBlocks() []*block {
	currentHash := GetBlockChain().NewestHash
	blocks := make([]*block, 0)
	for {
		if currentHash == "" {
			break
		}
		blocks = append(blocks, FindBlock(currentHash))
		currentHash = blocks[len(blocks)-1].PrevHash
	}
	return blocks
}
