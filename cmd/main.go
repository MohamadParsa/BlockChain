package main

import (
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/blocks/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/MohamadParsa/BlockChain/v1/miner/miner_server"
	"github.com/MohamadParsa/BlockChain/v1/wallet"
	"github.com/MohamadParsa/BlockChain/v1/wallet/wallet_server"
)

func main() {

	blockCh := blockchain.New()

	minerWallet, _ := wallet.New()
	miner := miner.New(blockCh, 3, minerWallet)
	minerServer := miner_server.New(miner)

	wA, _ := wallet.New()
	walletServerA := wallet_server.New(wA)
	wB, _ := wallet.New()
	walletServerB := wallet_server.New(wB)

	tr1, _ := wA.SendCrypto(wB.Address(), 1.1)
	sign1, _ := wA.Sign(tr1)
	result, err := blockCh.AddTransaction(wA.PublicKey(), sign1, tr1)
	fmt.Println(result, err)
	err = miner.Mining()

	blockCh.Print()
	fmt.Println(blockCh.CalculateTotalAmount(wA.Address()))
	fmt.Println(blockCh.CalculateTotalAmount(wB.Address()))
	fmt.Println(blockCh.CalculateTotalAmount(minerWallet.Address()))
	go minerServer.Serve(":8080")
	go walletServerA.Serve(":8090")
	walletServerB.Serve(":8091")

}
