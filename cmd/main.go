package main

import (
	"github.com/MohamadParsa/BlockChain/v1/blockchain"
)

func main() {
	blockCh := blockchain.New()
	blockCh.AddBlock(1)
	blockCh.AddBlock(2)
	blockCh.AddBlock(3)
	blockCh.Print()
}
