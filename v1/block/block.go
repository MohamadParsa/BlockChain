package block

import "time"

type Block struct {
	timestamp    int64
	nonce        int64
	previousHash string
	transactions []string
}

func New(nonce int64, previousHash string) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
	}
}
