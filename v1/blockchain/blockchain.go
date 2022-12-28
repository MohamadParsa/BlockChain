package blockchain

import (
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/block"
)

type BlockChain struct {
	transactionsPool []string
	chain            []*block.Block
}

func New() *BlockChain {
	blockChain := &BlockChain{}
	blockChain.AddBlock(0)
	return blockChain
}
func (blockChain *BlockChain) AddBlock(nonce int64) *block.Block {
	b := block.New(nonce, blockChain.LastBlock().Hash())
	blockChain.chain = append(blockChain.chain, b)
	return b
}

func (blockChain *BlockChain) Print() {
	for index, block := range blockChain.chain {
		fmt.Println("block number: ", index)
		block.Print()
	}
}
func (blockChain *BlockChain) LastBlock() *block.Block {
	if len(blockChain.chain) == 0 {
		return &block.Block{}
	}
	return blockChain.chain[len(blockChain.chain)-1]
}
