package main

import (
	"fmt"
	"github.com/google/jsonapi"
	"log"
	"net/http"
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
	Id   int    `jsonapi:primary,id"`
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

func main() {

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
			payload := new(Payload)

			if err := jsonapi.UnmarshalPayload(r.Body, payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			chain.AddBlock(*payload)
			jsonapi.MarshalPayload(w, &chain)

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
