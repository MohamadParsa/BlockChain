package main

import (
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

func main() {
	blockCh := blockchain.New()
	miner := miner.New(blockCh, 3, "miner_Address")

	blockCh.AddTransaction(transaction.New("sender1", "recipient1", 1.1))
	blockCh.AddTransaction(transaction.New("sender2", "recipient2", 2.1))
	blockCh.AddTransaction(transaction.New("sender3", "recipient3", 3.1))
	miner.Mining()
	blockCh.AddTransaction(transaction.New("sender1", "recipient2", 4.1))

	miner.Mining()
	blockCh.Print()
	fmt.Println(blockCh.CalculateTotalAmount("recipient2"))
}
