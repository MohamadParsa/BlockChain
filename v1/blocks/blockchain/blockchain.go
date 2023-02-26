package blockchain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/blocks/block"
	"github.com/MohamadParsa/BlockChain/v1/signature"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

type BlockChain struct {
	transactionsPool []*transaction.Transaction
	chain            []*block.Block
}

func New() *BlockChain {
	blockChain := &BlockChain{}
	blockChain.AddBlock(0)
	return blockChain
}
func (blockChain *BlockChain) AddBlock(nonce int64) *block.Block {
	b := block.New(nonce, blockChain.LastBlock().Hash(), blockChain.transactionsPool)
	blockChain.transactionsPool = []*transaction.Transaction{}
	blockChain.chain = append(blockChain.chain, b)
	return b
}

func (blockChain *BlockChain) Print() {
	for index, block := range blockChain.chain {
		fmt.Println("block number: ", index+1)
		block.Print()
	}
}
func (blockChain *BlockChain) LastBlock() *block.Block {
	if len(blockChain.chain) == 0 {
		return &block.Block{}
	}
	return blockChain.chain[len(blockChain.chain)-1]
}
func (blockChain *BlockChain) AddTransaction(publicKey *ecdsa.PublicKey, sign *signature.Signature, transaction *transaction.Transaction) (bool, error) {

	if ok, err := signature.VerifySignature(publicKey, sign, transaction); !ok {
		return ok, err
	}
	blockChain.transactionsPool = append(blockChain.transactionsPool, transaction)
	return true, nil
}
func (blockChain *BlockChain) CopyTransactions() []*transaction.Transaction {
	transactions := []*transaction.Transaction{}
	for _, transaction := range blockChain.transactionsPool {
		temp := *transaction
		transactions = append(transactions, &temp)
	}
	return transactions
}

func (blockChain *BlockChain) CalculateTotalAmount(walletAddress string) float64 {
	var totalAmount float64 = 0.0
	for _, block := range blockChain.chain {
		totalAmount = totalAmount + block.CalculateTotalAmount(walletAddress)
	}

	return totalAmount
}
