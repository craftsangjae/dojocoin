package cryto

import (
	"encoding/json"
	"errors"
	"github.com/craftsangjae/dojocoin/db"
	"github.com/craftsangjae/dojocoin/utils"
	strings "strings"
	"sync"
	"time"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5 // 난이도 재조정 주기
	blockInterval      int = 2 // 블록 별 시간
	allowedRange       int = 2
)

var chain *blockChain
var once sync.Once

func init() {
	GetBlockChain()
}

type block struct {
	Transactions []*Tx
	Hash         string
	PrevHash     string
	Height       int
	Difficulty   int
	Nonce        int
	Timestamp    int
}

type blockChain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
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

func (chain *blockChain) recalculateDifficulty() int {
	var err error
	newstBlock, err := findLatestBlock(0)
	utils.HandleErr(err)
	lastRecalculatedBlock, err := findLatestBlock(difficultyInterval - 1)
	utils.HandleErr(err)
	actualTime := (newstBlock.Timestamp - lastRecalculatedBlock.Timestamp) / 60
	expectedTime := difficultyInterval * blockInterval

	if actualTime < expectedTime-allowedRange {
		return chain.CurrentDifficulty + 1
	} else if actualTime > expectedTime+allowedRange {
		return chain.CurrentDifficulty - 1
	}
	return chain.CurrentDifficulty
}

// 난이도를 계산합니다

func (chain *blockChain) difficulty() int {
	if chain.Height == 0 {
		return defaultDifficulty
	} else if chain.Height%5 == 0 {
		// recalculate the difficulty
		return chain.recalculateDifficulty()
	} else {
		return chain.CurrentDifficulty
	}
}

// 새 블록을 생성합니다
//
func newBlock() *block {
	prevHash := chain.NewestHash
	return &block{nil, "", prevHash, chain.Height + 1, chain.difficulty(), 0, 0}
}

func initBlock() *block {
	return &block{[]*Tx{}, "", "", 0, 0, 0, 0}
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
	chain.CurrentDifficulty = newBlock.Difficulty
}

func UTxOutsByAddress(address string) (utxOuts []*UTxOut) {
	creatorTxs := make(map[string]bool)
	for _, block := range ListAllBlocks() {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						utxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(utxOut) {
							utxOuts = append(utxOuts, utxOut)
						}
					}
				}
			}
		}
	}
	return
}

func BalanceByAddress(address string) (output int) {
	for _, tx := range UTxOutsByAddress(address) {
		output += tx.Amount
	}
	return
}

func AddBlock() {
	b := newBlock()
	b.mine()
	b.Transactions = Mempool.TxToConfirm()
	b.persist()

	chain.update(b)
	chain.persist()
}

func FindBlock(hash string) *block {
	b := initBlock()
	data := db.FindData(hash)
	utils.FromBytes(b, data)
	return b
}

func GetBlockChain() *blockChain {
	if chain == nil {
		once.Do(func() {
			checkpoint := db.LoadCheckpoint()
			chain = &blockChain{"", 0, 0}
			if checkpoint == nil {
				AddBlock()
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

func findLatestBlock(height int) (*block, error) {
	return findBlockRecursive(height, GetBlockChain().NewestHash)
}

func findBlockRecursive(height int, hash string) (*block, error) {
	b := FindBlock(hash)
	if b == nil {
		return nil, errors.New("Not Found")
	}

	if height == 0 {
		return b, nil
	} else if height > 0 {
		return findBlockRecursive(height-1, b.PrevHash)
	} else {
		return nil, errors.New("Error")
	}
}
