package zombie

import (
	"../chain"
)

type Zombie struct {
	Node *chain.Node
}

var zpool []*Zombie

func generateIP() string {
	return "1.1.1.1:3001"
}

func generateUsername() string {
	return "Cpt. Sr. McGilligan III"
}

func generateZombies(n int) {
	for i := 1; i <= n; i++ {
		z := &Zombie{
			Node: &chain.Node{
				IP:       generateIP(),
				Username: generateUsername(),
			},
		}

		zpool = append(zpool, z)
	}
}

func main() {
	nOfZs := 3 //grab this from the CMD eventually
	generateZombies(nOfZs)
}
