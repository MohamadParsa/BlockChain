package transaction_request

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

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
func (transactionRequest *TransactionRequest) PublicKeyString() string {
	return fmt.Sprintf("%x%x", transactionRequest.PublicKey.X, transactionRequest.PublicKey.Y)
}
func (transactionRequest *TransactionRequest) Print() {
	fmt.Printf("$	sender:		%s\n", transactionRequest.SenderAddress())
	fmt.Printf("$	recipient:	%s\n", transactionRequest.RecipientAddress())
	fmt.Printf("$	value:		%3f\n", transactionRequest.Value())
	fmt.Printf("$	signature:		%3s\n", transactionRequest.Signature)
	fmt.Printf("$	publicKey:		%3s\n", transactionRequest.PublicKey)
}
func (transactionRequest *TransactionRequest) UnMarshalJSON() {
}

func (transactionRequest *TransactionRequest) decodePublicKeyDataFromString(signatureText string) (*big.Int, *big.Int, error) {
	var x, y big.Int
	byteX, err := hex.DecodeString(signatureText[:64])
	if err != nil {
		return nil, nil, err
	}
	byteY, err := hex.DecodeString(signatureText[64:])
	if err != nil {
		return nil, nil, err
	}

	return x.SetBytes(byteX), y.SetBytes(byteY), nil
}
