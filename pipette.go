package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func (block *Block) save() error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Block"))

		if block.Index <= 0 {
			id, _ := b.NextSequence()
			block.Index = int(id)
		}

		encoded, err := json.Marshal(block)
		if err != nil {
			return err
		}

		err = b.Put([]byte(string(block.Index)), encoded)
		return err
	})
}

func (n *Node) save() error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Node"))

		if n.ID <= 0 {
			id, _ := b.NextSequence()
			n.ID = int(id)
		}

		encoded, err := json.Marshal(n)
		if err != nil {
			return err
		}

		err = b.Put([]byte(string(n.ID)), encoded)
		return err
	})
}

func loadChain() Chain {
	blocks := []*Block{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Block"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var block Block
			json.Unmarshal(v, &block)
			blocks = append(blocks, &block)
		}

		return nil
	})

	chain := Chain{
		Blocks: blocks,
	}

	return chain
}

func createBuckets() error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Block"))
		_, err = tx.CreateBucketIfNotExists([]byte("Node"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func noBlocksExist() bool {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Block"))
		v := b.Get([]byte(string(1)))

		if v == nil {
			return errors.New("No blocks found")
		}

		return nil
	})

	return err != nil
}

func setupDB() {
	tempDB, err := bolt.Open("my.db", 0600, nil)
	err = createBuckets()
	if err != nil {
		log.Fatal(err)
	}

	db = tempDB
}
