package main

import (
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/blocks/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/MohamadParsa/BlockChain/v1/miner/server"
	"github.com/MohamadParsa/BlockChain/v1/wallet"
)

func main() {

	blockCh := blockchain.New()

	minerWallet, _ := wallet.New()
	miner := miner.New(blockCh, 3, minerWallet)
	restFull := server.New(miner)

	wA, _ := wallet.New()
	wB, _ := wallet.New()

	tr1, _ := wA.SendCrypto(wB.Address(), 1.1)
	sign1, _ := wA.Sign(tr1)
	result, err := blockCh.AddTransaction(wA.PublicKey(), sign1, tr1)
	fmt.Println(result, err)
	err = miner.Mining()

	blockCh.Print()
	fmt.Println(blockCh.CalculateTotalAmount(wA.Address()))
	fmt.Println(blockCh.CalculateTotalAmount(wB.Address()))
	fmt.Println(blockCh.CalculateTotalAmount(minerWallet.Address()))
	restFull.Serve(":8080")
}
