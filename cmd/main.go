package main

import (
	"github.com/MohamadParsa/BlockChain/v1/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

func main() {
	blockCh := blockchain.New()
	miner := miner.New(blockCh, 3)
	blockCh.AddTransaction(transaction.New("sender1", "recipient1", 1.1))
	blockCh.AddTransaction(transaction.New("sender2", "recipient2", 2.1))
	blockCh.AddTransaction(transaction.New("sender3", "recipient3", 3.1))
	nonce := miner.FindValidNonce()
	blockCh.AddBlock(nonce)
	nonce = miner.FindValidNonce()
	blockCh.AddBlock(nonce)
	blockCh.Print()
}
