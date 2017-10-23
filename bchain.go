package main

import (
	"github.com/google/jsonapi"
	"log"
	"net/http"
)

type Payload struct {
	Data string `jsonapi:"attr, data"`
}

type Block struct {
	Index    int `jsonapi:"primary,block"`
	Payload  `jsonapi:"relation,payload"`
	Previous *Block `jsonapi:"relation,previous"`
}

func main() {

	block := Block{Index: 0, Payload: Payload{Data: "This is my payload"}, Previous: nil}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		//w.WriteHeader(http.StatusOK)
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		//fmt.Fprintf(w, "Hello, %q", block.data)
		if err := jsonapi.MarshalPayload(w, &block); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
