package main

import (
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/blocks/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
	"github.com/MohamadParsa/BlockChain/v1/wallet"
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

	w, err := wallet.New()
	fmt.Println(w, err, w.PublicKey())
	transaction, _ := w.SendCrypto("a", 1)
	fmt.Println(transaction)
	fmt.Println("ffffffff")
	fmt.Println(w.Sign(transaction))
	fmt.Println("w.Sign(transaction)")
	fmt.Println(w.Address())

}
