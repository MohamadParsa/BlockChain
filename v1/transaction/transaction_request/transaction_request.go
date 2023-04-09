package transaction_request

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/MohamadParsa/BlockChain/v1/signature"
	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

type TransactionRequest struct {
	transaction.Transaction
	signature *signature.Signature `json:"-"`
	publicKey *ecdsa.PublicKey     `json:"-"`
	Signature string               `json:"signature"`
	PublicKey string               `json:"publicKey"`
}

func NewTransactionRequest(transaction transaction.Transaction, signature *signature.Signature, publicKey *ecdsa.PublicKey) *TransactionRequest {
	transactionRequest := &TransactionRequest{
		Transaction: transaction,
		signature:   signature,
		publicKey:   publicKey,
	}
	transactionRequest.PublicKey = transactionRequest.PublicKeyString()
	transactionRequest.Signature = transactionRequest.SignatureString()
	return transactionRequest
}
func (transactionRequest *TransactionRequest) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	type TransactionRequestTemp struct {
		transaction.Transaction
		Signature string `json:"signature"`
		PublicKey string `json:"publicKey"`
	}
	transactionRequestTemp := &TransactionRequestTemp{}
	err := json.Unmarshal(data, transactionRequestTemp)
	if err != nil {
		return err
	}
	signature, err := transactionRequest.signature.DecodeSignature(transactionRequestTemp.Signature)
	if err != nil {
		return err
	}
	publicKey, err := transactionRequest.decodePublicKey(transactionRequestTemp.PublicKey)
	if err != nil {
		return err
	}
	transactionRequest.Transaction = transactionRequestTemp.Transaction
	transactionRequest.signature = signature
	transactionRequest.publicKey = publicKey
	transactionRequest.Signature = signature.String()
	transactionRequest.PublicKey = fmt.Sprintf("%x%x", publicKey.X, publicKey.Y)
	return err
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
	return fmt.Sprintf("%x%x", transactionRequest.publicKey.X, transactionRequest.publicKey.Y)
}
func (transactionRequest *TransactionRequest) SignatureString() string {
	return transactionRequest.signature.String()
}
func (transactionRequest *TransactionRequest) GetSignature() signature.Signature {
	return *transactionRequest.signature
}
func (transactionRequest *TransactionRequest) GetPublicKey() ecdsa.PublicKey {
	return *transactionRequest.publicKey
}
func (transactionRequest *TransactionRequest) Print() {
	fmt.Printf("$	sender:		%s\n", transactionRequest.SenderAddress())
	fmt.Printf("$	recipient:	%s\n", transactionRequest.RecipientAddress())
	fmt.Printf("$	value:		%3f\n", transactionRequest.Value())
	fmt.Printf("$	signature:		%3s\n", transactionRequest.Signature)
	fmt.Printf("$	publicKey:		%3s\n", transactionRequest.PublicKey)
}

func (transactionRequest *TransactionRequest) decodePublicKey(publicKeyString string) (*ecdsa.PublicKey, error) {
	var x, y big.Int
	byteX, err := hex.DecodeString(publicKeyString[:64])
	if err != nil {
		return nil, err
	}
	byteY, err := hex.DecodeString(publicKeyString[64:])
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{elliptic.P256(), x.SetBytes(byteX), y.SetBytes(byteY)}, nil
}
