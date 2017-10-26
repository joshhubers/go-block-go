package main

import (
	"fmt"
	"github.com/google/jsonapi"
	"log"
	"net"
	"net/http"
	"os"
	"bufio"
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
	ID   int    `jsonapi:"primary,payload"`
	Info string `jsonapi:"attr,info"`
}

type Node struct {
	ID int `jsonapi:"primary,node"`
	IP int `jsonapi:"attr,ip"`
}

type NodeList struct {
	Nodes []*Node `jsonapi:"relation,nodes"`
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

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}


	myIPs := []string{}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				myIPs = append(myIPs, ipnet.IP.String())
			}
		}
	}

	if len(myIPs) == 0 {
		os.Stderr.WriteString("You are not connected to a network.\n")
		os.Exit(1)
	}

	return myIPs[0]
}

func determineStartState() string {
	fmt.Println("Enter the IP address you wish to connect to, or leave blank to start your own chain")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	if text == "\n" {
		return ""
	} else {
		return text
	}
}

func generateGenesis() Chain {
	block := Block{
		Index: 0,
		Data: &Payload{
			ID:   0,
			Info: "This is my payload",
		},
		Previous: nil,
	}

	chain := Chain{
		Blocks: []*Block{&block},
	}

	return chain
}

func beginServer(chain Chain) {
	http.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)

		if r.Method == "POST" {
			previousBlock := chain.Blocks[len(chain.Blocks)-1]

			payload := Payload{
				ID: previousBlock.Data.ID + 1,
			}

			if err := jsonapi.UnmarshalPayload(r.Body, &payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(payload)
			chain.AddBlock(payload)
			jsonapi.MarshalPayload(w, &chain)

		}
	})

	http.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		if err := jsonapi.MarshalPayload(w, &chain); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		if r.Method == "POST" {
			node := Node{
				ID: 1,
			}

			if err := jsonapi.UnmarshalPayload(r.Body, &node); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(node)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func goGetChain(targetIP string) Chain {
	fmt.Println("Going to get Chain from " + targetIP)
	chain := Chain{}
	return chain
}

func main() {
	myIP := getIP()
	targetIP := determineStartState()

	chain := Chain{}

	if len(targetIP) == 0 {
		chain = generateGenesis()
	} else {
		chain = goGetChain(targetIP)
	}

	fmt.Println("Your IP address for others to connect to is: " + myIP)
	beginServer(chain)
}
