package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/google/jsonapi"
	"golang.org/x/crypto/bcrypt"
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
	Info string `jsonapi:"attr,info"`
}

type Node struct {
	ID       int    `jsonapi:"primary,node"`
	IP       int    `jsonapi:"attr,ip"`
	Username string `jsonapi:"attr,username"`
	Hash     []byte `jsonapi:"attr, hash"`
}

type NodeList struct {
	Nodes []*Node `jsonapi:"relation,nodes"`
}

func (bchain *Chain) addBlock(data Payload) []*Block {
	previousBlock := bchain.Blocks[len(bchain.Blocks)-1]
	newBlock := &Block{
		Index:    previousBlock.Index + 1,
		Data:     &data,
		Previous: previousBlock,
	}

	newBlock.save()

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
	}

	return text
}

func generateGenesis(db *bolt.DB) Chain {
	chain := Chain{}
	if noBlocksExist() {
		fmt.Println("No genesis block found, creating new one...")
		block := Block{
			Index: 0,
			Data: &Payload{
				Info: "In the beginning...",
			},
			Previous: nil,
		}

		block.save()

		chain = Chain{
			Blocks: []*Block{&block},
		}
	} else {
		fmt.Println("Previous blocks detected loading in blocks...")
		chain = loadChain()
	}

	return chain
}

func beginServer(chain Chain) {
	http.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonapi.MediaType)

		if r.Method == "POST" {
			payload := Payload{}

			if err := jsonapi.UnmarshalPayload(r.Body, &payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(payload)
			chain.addBlock(payload)
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
			node := Node{}

			if err := jsonapi.UnmarshalPayload(r.Body, &node); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(node.Username) > 0 {
				http.Error(w, "Username must be provided", http.StatusBadRequest)
			}

			//TODO: Check if username-ip is used already
			//TODO: Check that username-ip is provided
			hash, herr := bcrypt.GenerateFromPassword([]byte(string(node.Username)+string(node.IP)), bcrypt.DefaultCost)
			if herr != nil {
				// TODO: Properly handle error
				log.Fatal(herr)
			}

			node.Hash = hash
			node.save()
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
	setupDB()
	defer db.Close()

	myIP := getIP()
	targetIP := determineStartState()

	chain := Chain{}

	if len(targetIP) == 0 {
		chain = generateGenesis(db)
	} else {
		chain = goGetChain(targetIP)
	}

	fmt.Println("Your IP address for others to connect to is: " + myIP)
	beginServer(chain)
}
