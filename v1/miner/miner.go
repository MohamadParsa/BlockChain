package miner

import (
	"fmt"
	"strings"

	"github.com/MohamadParsa/BlockChain/v1/blocks/block"
	"github.com/MohamadParsa/BlockChain/v1/blocks/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

const (
	MINING_REWARD_SENDER_ADDRESS = "reward_sender_address"
	REWARD                       = 6.25
)

type Miner struct {
	blockChain         *blockchain.BlockChain
	difficulty         int
	minerWalletAddress string
}

func New(blockChain *blockchain.BlockChain, difficulty int, minerWalletAddress string) *Miner {
	return &Miner{
		blockChain:         blockChain,
		difficulty:         difficulty,
		minerWalletAddress: minerWalletAddress,
	}
}

func (miner *Miner) Mining() {

	transaction := transaction.New(MINING_REWARD_SENDER_ADDRESS, miner.minerWalletAddress, REWARD)
	miner.blockChain.AddTransaction(transaction)

	transactions := miner.blockChain.CopyTransactions()
	previousHash := miner.blockChain.LastBlock().Hash()

	nonce := miner.findValidNonce(previousHash, transactions)

	miner.blockChain.AddBlock(nonce)

}

func (miner *Miner) findValidNonce(previousHash [32]byte, transactions []*transaction.Transaction) int64 {
	var nonce int64 = 1
	findNumber := make(chan int64)

	for {
		select {
		case v, ok := <-findNumber:
			if ok {
				return v
			}
		default:
			go miner.validateNonce(nonce, previousHash, transactions, findNumber)
			nonce++
		}
	}
}

func (miner *Miner) validateNonce(nonce int64, previousHash [32]byte, transactions []*transaction.Transaction, findNumber chan int64) {
	zeros := strings.Repeat("0", miner.difficulty)

	candidateBlock := block.CreateCandidateBlock(nonce, previousHash, transactions)
	if candidateBlock.HashString()[:miner.difficulty] == zeros {
		fmt.Println(candidateBlock.HashString())
		findNumber <- nonce
	}

}
