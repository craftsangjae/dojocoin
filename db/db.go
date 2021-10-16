package db

import (
	"github.com/boltdb/bolt"
	. "github.com/craftsangjae/dojocoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

var db *bolt.DB

func DB() (dbPointer *bolt.DB) {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
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
		return bucket.Put([]byte(hash), ToBytes(data))
	})
	HandleErr(err)
}

func SaveBlockchain(data interface{}) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		return bucket.Put([]byte("checkpoint"), ToBytes(data))
	})
	HandleErr(err)
}
