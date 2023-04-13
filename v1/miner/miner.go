package miner

import (
	"fmt"
	"strings"

	"github.com/MohamadParsa/BlockChain/v1/blocks/block"
	"github.com/MohamadParsa/BlockChain/v1/blocks/blockchain"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
	transaction_request "github.com/MohamadParsa/BlockChain/v1/transaction/transactionDTO"
	"github.com/MohamadParsa/BlockChain/v1/wallet"
)

type Miner struct {
	blockChain  *blockchain.BlockChain
	difficulty  int
	minerWallet *wallet.Wallet
}

func New(blockChain *blockchain.BlockChain, difficulty int, minerWallet *wallet.Wallet) *Miner {

	return &Miner{
		blockChain:  blockChain,
		difficulty:  difficulty,
		minerWallet: minerWallet,
	}
}

func (miner *Miner) Mining() error {
	if miner.blockChain.TransactionPoolCount() == 0 {
		return nil
	}
	transaction := transaction.New(miner.blockChain.MiningRewardSenderAddress(), miner.minerWallet.Address(), miner.blockChain.MiningReward())

	miner.blockChain.AddTransaction(miner.minerWallet.PublicKey(), nil, transaction)

	transactions := miner.blockChain.CopyTransactions()

	previousHash := miner.blockChain.LastBlock().Hash()

	nonce := miner.findValidNonce(previousHash, transactions)

	miner.blockChain.AddBlock(nonce)
	return nil

}
func (miner *Miner) AddTransaction(transactionDTO *transaction_request.TransactionDTO) (bool, error) {
	publicKey := transactionDTO.GetPublicKey()
	signature := transactionDTO.GetSignature()
	ok, err := miner.blockChain.AddTransaction(&publicKey, &signature, &transactionDTO.Transaction)
	return ok, err

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

func (miner *Miner) BlockChain() *blockchain.BlockChain {
	return miner.blockChain

}
