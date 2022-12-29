package miner

import (
	"fmt"
	"strings"

	"github.com/MohamadParsa/BlockChain/v1/block"
	"github.com/MohamadParsa/BlockChain/v1/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

type Miner struct {
	blockChain *blockchain.BlockChain
	difficulty int
}

func New(blockChain *blockchain.BlockChain, difficulty int) *Miner {
	return &Miner{
		blockChain: blockChain,
		difficulty: difficulty,
	}
}

func (miner *Miner) FindValidNonce() int64 {

	transactions := miner.blockChain.CopyTransactions()
	previousHash := miner.blockChain.LastBlock().Hash()

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

	candidateBlock := block.New(nonce, previousHash, transactions)
	if candidateBlock.HashString()[:miner.difficulty] == zeros {
		fmt.Println(candidateBlock.HashString())
		findNumber <- nonce
	}

}
