package transaction

import (
	"encoding/json"
	"fmt"
)

type TransactionRequest struct {
	Transaction
	signature string
	publicKey string
}

func NewTransactionRequest(transaction Transaction, signature string, publicKey string) *TransactionRequest {
	return &TransactionRequest{
		Transaction: transaction,
		signature:   signature,
		publicKey:   publicKey,
	}
}

func (transactionRequest *TransactionRequest) SenderAddress() string {
	return transactionRequest.senderAddress
}

func (transactionRequest *TransactionRequest) RecipientAddress() string {
	return transactionRequest.recipientAddress
}

func (transactionRequest *TransactionRequest) Value() float64 {
	return transactionRequest.value
}
func (transactionRequest *TransactionRequest) Signature() string {
	return transactionRequest.signature
}
func (transactionRequest *TransactionRequest) PublicKey() string {
	return transactionRequest.publicKey
}
func (transactionRequest *TransactionRequest) Print() {
	fmt.Printf("$	sender:		%s\n", transactionRequest.senderAddress)
	fmt.Printf("$	recipient:	%s\n", transactionRequest.recipientAddress)
	fmt.Printf("$	value:		%3f\n", transactionRequest.value)
	fmt.Printf("$	signature:		%3s\n", transactionRequest.signature)
	fmt.Printf("$	publicKey:		%3s\n", transactionRequest.publicKey)
}
func (transactionRequest *TransactionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float64 `json:"value"`
		Signature string  `json:"signature"`
		PublicKey string  `json:"publicKey"`
	}{
		Sender:    transactionRequest.senderAddress,
		Recipient: transactionRequest.recipientAddress,
		Value:     transactionRequest.value,
		Signature: transactionRequest.signature,
		PublicKey: transactionRequest.publicKey,
	})
}
