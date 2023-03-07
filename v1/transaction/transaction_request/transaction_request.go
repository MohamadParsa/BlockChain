package transaction_request

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/MohamadParsa/BlockChain/v1/signature"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

type TransactionRequest struct {
	transaction.Transaction
	Signature *signature.Signature `json:"signature"`
	PublicKey *ecdsa.PublicKey     `json:"publicKey"`
}

func NewTransactionRequest(transaction transaction.Transaction, signature *signature.Signature, publicKey *ecdsa.PublicKey) *TransactionRequest {
	return &TransactionRequest{
		Transaction: transaction,
		Signature:   signature,
		PublicKey:   publicKey,
	}
}

func (transactionRequest *TransactionRequest) SenderAddress() string {
	return transactionRequest.Transaction.SenderAddress
}

func (transactionRequest *TransactionRequest) RecipientAddress() string {
	return transactionRequest.Transaction.RecipientAddress
}

func (transactionRequest *TransactionRequest) Value() float64 {
	return transactionRequest.Transaction.Value
}

func (transactionRequest *TransactionRequest) Print() {
	fmt.Printf("$	sender:		%s\n", transactionRequest.SenderAddress())
	fmt.Printf("$	recipient:	%s\n", transactionRequest.RecipientAddress())
	fmt.Printf("$	value:		%3f\n", transactionRequest.Value())
	fmt.Printf("$	signature:		%3s\n", transactionRequest.Signature)
	fmt.Printf("$	publicKey:		%3s\n", transactionRequest.PublicKey)
}
func (transactionRequest *TransactionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string               `json:"senderAddress"`
		Recipient string               `json:"recipientAddress"`
		Value     float64              `json:"value"`
		Signature *signature.Signature `json:"signature"`
		PublicKey *ecdsa.PublicKey     `json:"publicKey"`
	}{
		Sender:    transactionRequest.SenderAddress(),
		Recipient: transactionRequest.RecipientAddress(),
		Value:     transactionRequest.Value(),
		Signature: transactionRequest.Signature,
		PublicKey: transactionRequest.PublicKey,
	})
}
