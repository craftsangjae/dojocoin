package db

import (
	"github.com/boltdb/bolt"
	. "github.com/craftsangjae/dojocoin/utils"
)

const (
	dbName           = "blockchain.db"
	checkpointBucket = "data"
	blocksBucket     = "blocks"
	checkpointKey    = "checkpoint"
)

var db *bolt.DB

func DB() (dbPointer *bolt.DB) {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(checkpointBucket))
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		HandleErr(err)
	}
	return db
}

func SaveBlock(hash string, data interface{}) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		dataBytes := ToBytes(data)
		return bucket.Put([]byte(hash), dataBytes)
	})
	HandleErr(err)
}

func SaveBlockchain(data interface{}) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(checkpointBucket))
		return bucket.Put([]byte("checkpoint"), ToBytes(data))
	})
	HandleErr(err)
}

func LoadCheckpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(checkpointBucket))
		data = bucket.Get([]byte(checkpointKey))
		return nil
	})
	return data
}

func FindData(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
