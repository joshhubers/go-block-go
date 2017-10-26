package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/google/jsonapi"
)

type Chain struct {
	Blocks []*Block `jsonapi:"relation,blocks"`
}

type Block struct {
	Index    int      `jsonapi:"primary,block"`
	Data     *Payload `jsonapi:"relation,data"`
	Previous *Block   `jsonapi:"relation,previous"`
}

type Payload struct {
	Id   int    `jsonapi:"primary,payload"`
	Info string `jsonapi:"attr,info"`
}

func (bchain *Chain) AddBlock(data Payload) []*Block {
	previousBlock := bchain.Blocks[len(bchain.Blocks)-1]
	newBlock := &Block{
		Index:    previousBlock.Index + 1,
		Data:     &data,
		Previous: previousBlock,
	}

	bchain.Blocks = append(bchain.Blocks, newBlock)
	return bchain.Blocks
}

func (block *Block) save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Block"))
		err := b.Put([]byte("foo"), []byte("42"))
		return err
	})
}

func createBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Block"))
		_, err = tx.CreateBucketIfNotExists([]byte("Chain"))
		_, err = tx.CreateBucketIfNotExists([]byte("Payload"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func setupDB() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	err = createBuckets(db)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := setupDB()

	defer db.Close()
	block := Block{
		Index: 0,
		Data: &Payload{
			Id:   0,
			Info: "This is my payload",
		},
		Previous: nil,
	}

	testP := Payload{
		Id:   1,
		Info: "My new info",
	}

	chain := Chain{
		Blocks: []*Block{&block},
	}

	chain.AddBlock(testP)

	fmt.Println(chain)

	http.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			previousBlock := chain.Blocks[len(chain.Blocks)-1]

			payload := Payload{
				Id: previousBlock.Data.Id + 1,
			}

			if err := jsonapi.UnmarshalPayload(r.Body, &payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(payload)
			chain.AddBlock(payload)
			jsonapi.MarshalPayload(w, &chain)

		}

		if r.Method == "GET" {
			jsonapi.MarshalPayload(w, &testP)
		}
	})

	http.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		if err := jsonapi.MarshalPayload(w, &chain); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
