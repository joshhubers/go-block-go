package main

import (
	"github.com/google/jsonapi"
	"log"
	"net/http"
)

type Block struct {
	Index    int      `jsonapi:"primary,block"`
	Data     *Payload `jsonapi:"relation,data"`
	Previous *Block   `jsonapi:"relation,previous"`
}

type Payload struct {
	Id   int    `jsonapi:primary,id"`
	Info string `jsonapi:"attr,info"`
}

func main() {

	block := &Block{
		Index: 0,
		Data: &Payload{
			Id:   0,
			Info: "This is my payload",
		},
		Previous: nil,
	}

	http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		if err := jsonapi.MarshalPayload(w, block); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
