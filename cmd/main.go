package main

import (
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

	go minerServer.Serve(":8080")
	go walletServerA.Serve(":8090")
	walletServerB.Serve(":8091")

}
