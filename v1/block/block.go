package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int64
	previousHash [32]byte
	transactions []string
}

func New(nonce int64, previousHash [32]byte) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
	}
}

func (block *Block) Print() {

	fmt.Println(strings.Repeat("_", 100))
	fmt.Printf("| timestamp:		%d\n", block.timestamp)
	fmt.Printf("| nonce:		%d\n", block.nonce)
	fmt.Printf("| previous hash:	%x\n", block.previousHash)
	fmt.Printf("| transactions:		%s\n", block.transactions)
	fmt.Println(strings.Repeat("-", 100))

}
func (block *Block) Hash() [32]byte {
	marshal, _ := json.Marshal(block)
	return sha256.Sum256(marshal)
}
func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        int64    `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Transactions []string `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: block.previousHash,
		Transactions: block.transactions,
	})
}
