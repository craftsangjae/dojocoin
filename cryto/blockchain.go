package cryto

import (
	"encoding/json"
	"github.com/craftsangjae/dojocoin/db"
	"github.com/craftsangjae/dojocoin/utils"
	strings "strings"
	"sync"
	"time"
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
	Timestamp  int
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
	return &block{data, "", prevHash, chain.Height + 1, difficulty, 0, 0}
}

func (b *block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		}
		b.Nonce++
	}
}

func (chain *blockChain) update(newBlock *block) {
	chain.NewestHash = newBlock.Hash
	chain.Height = newBlock.Height
}

func AddBlock(data string) {
	b := newBlock(data)
	b.mine()
	b.persist()

	chain.update(b)
	chain.persist()
}

func FindBlock(hash string) *block {
	b := &block{"", "", "", 1, difficulty, 0, 0}
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
