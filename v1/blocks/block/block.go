package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MohamadParsa/BlockChain/v1/blocks/transaction"
)

type Block struct {
	timestamp    int64
	nonce        int64
	previousHash [32]byte
	transactions []*transaction.Transaction
}

func New(nonce int64, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
}
func CreateCandidateBlock(nonce int64, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
}
func (block *Block) Print() {

	fmt.Println(strings.Repeat("_", 100))
	fmt.Printf("| timestamp:		%d\n", block.timestamp)
	fmt.Printf("| nonce:		%d\n", block.nonce)
	fmt.Printf("| previous hash:	%x\n", block.previousHash)
	for index, transaction := range block.transactions {
		fmt.Printf("| transactions number:		%d\n", index+1)
		transaction.Print()
	}
	fmt.Println(strings.Repeat("-", 100))

}
func (block *Block) HashString() string {
	hash := block.Hash()
	return fmt.Sprintf("%x", hash)
}
func (block *Block) Hash() [32]byte {
	marshal, _ := json.Marshal(block)
	return sha256.Sum256(marshal)
}
func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                      `json:"timestamp"`
		Nonce        int64                      `json:"nonce"`
		PreviousHash [32]byte                   `json:"previous_hash"`
		Transactions []*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: block.previousHash,
		Transactions: block.transactions,
	})
}
func (block *Block) CalculateTotalAmount(walletAddress string) float64 {
	var totalAmount float64 = 0.0
	for _, transaction := range block.transactions {
		if transaction.SenderAddress() == walletAddress {
			totalAmount = totalAmount - transaction.Value()
		}
		if transaction.RecipientAddress() == walletAddress {
			totalAmount = totalAmount + transaction.Value()
		}
	}
	return totalAmount
}
