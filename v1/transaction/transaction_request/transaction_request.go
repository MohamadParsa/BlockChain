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
	signature *signature.Signature
	publicKey *ecdsa.PublicKey
}

func NewTransactionRequest(transaction transaction.Transaction, signature *signature.Signature, publicKey *ecdsa.PublicKey) *TransactionRequest {
	return &TransactionRequest{
		Transaction: transaction,
		signature:   signature,
		publicKey:   publicKey,
	}
}

func (transactionRequest *TransactionRequest) SenderAddress() string {
	return transactionRequest.SenderAddress()
}

func (transactionRequest *TransactionRequest) RecipientAddress() string {
	return transactionRequest.Transaction.RecipientAddress()
}

func (transactionRequest *TransactionRequest) Value() float64 {
	return transactionRequest.Value()
}
func (transactionRequest *TransactionRequest) Signature() signature.Signature {
	return *transactionRequest.signature
}
func (transactionRequest *TransactionRequest) PublicKey() ecdsa.PublicKey {
	return *transactionRequest.publicKey
}
func (transactionRequest *TransactionRequest) Print() {
	fmt.Printf("$	sender:		%s\n", transactionRequest.SenderAddress())
	fmt.Printf("$	recipient:	%s\n", transactionRequest.RecipientAddress())
	fmt.Printf("$	value:		%3f\n", transactionRequest.Value())
	fmt.Printf("$	signature:		%3s\n", transactionRequest.signature)
	fmt.Printf("$	publicKey:		%3s\n", transactionRequest.publicKey)
}
func (transactionRequest *TransactionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string               `json:"sender_address"`
		Recipient string               `json:"recipient_address"`
		Value     float64              `json:"value"`
		Signature *signature.Signature `json:"signature"`
		PublicKey *ecdsa.PublicKey     `json:"publicKey"`
	}{
		Sender:    transactionRequest.SenderAddress(),
		Recipient: transactionRequest.RecipientAddress(),
		Value:     transactionRequest.Value(),
		Signature: transactionRequest.signature,
		PublicKey: transactionRequest.publicKey,
	})
}
