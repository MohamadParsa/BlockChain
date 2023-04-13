package transactionDTO

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

type TransactionDTO struct {
	transaction.Transaction
	signature *signature.Signature `json:"-"`
	publicKey *ecdsa.PublicKey     `json:"-"`
	Signature [2]string            `json:"signature"`
	PublicKey [2]string            `json:"publicKey"`
}

func NewTransactionDTO(transaction transaction.Transaction, signature *signature.Signature, publicKey *ecdsa.PublicKey) *TransactionDTO {
	transactionDTO := &TransactionDTO{
		Transaction: transaction,
		signature:   signature,
		publicKey:   publicKey,
	}
	transactionDTO.PublicKey = transactionDTO.PublicKeyString()
	transactionDTO.Signature = transactionDTO.SignatureString()
	return transactionDTO
}
func (transactionDTO *TransactionDTO) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	type TransactionDTOTemp struct {
		Sender    string    `json:"senderAddress"`
		Recipient string    `json:"recipientAddress"`
		Value     float64   `json:"value"`
		Signature [2]string `json:"signature"`
		PublicKey [2]string `json:"publicKey"`
	}
	transactionDTOTemp := &TransactionDTOTemp{}
	err := json.Unmarshal(data, transactionDTOTemp)
	if err != nil {
		return err
	}

	transaction := transaction.New(transactionDTOTemp.Sender, transactionDTOTemp.Recipient, transactionDTOTemp.Value)

	signature, err := transactionDTO.signature.DecodeSignature(transactionDTOTemp.Signature)
	if err != nil {
		return err
	}

	publicKey, err := transactionDTO.DecodePublicKey(transactionDTOTemp.PublicKey)
	if err != nil {
		return err
	}
	transactionDTO.Transaction = *transaction
	transactionDTO.signature = signature
	transactionDTO.publicKey = publicKey
	transactionDTO.Signature = signature.String()
	transactionDTO.PublicKey = [2]string{fmt.Sprintf("%x", publicKey.X), fmt.Sprintf("%x", publicKey.Y)}

	return err
}
func (transactionDTO *TransactionDTO) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string    `json:"senderAddress"`
		Recipient string    `json:"recipientAddress"`
		Value     float64   `json:"value"`
		Signature [2]string `json:"signature"`
		PublicKey [2]string `json:"publicKey"`
	}{
		Sender:    transactionDTO.Transaction.SenderAddress,
		Recipient: transactionDTO.Transaction.RecipientAddress,
		Value:     transactionDTO.Transaction.Value,
		Signature: transactionDTO.SignatureString(),
		PublicKey: transactionDTO.PublicKeyString(),
	})
}
func (transactionDTO *TransactionDTO) SenderAddress() string {
	return transactionDTO.Transaction.SenderAddress
}

func (transactionDTO *TransactionDTO) RecipientAddress() string {
	return transactionDTO.Transaction.RecipientAddress
}

func (transactionDTO *TransactionDTO) Value() float64 {
	return transactionDTO.Transaction.Value
}
func (transactionDTO *TransactionDTO) PublicKeyString() [2]string {
	return [2]string{fmt.Sprintf("%x", transactionDTO.publicKey.X), fmt.Sprintf("%x", transactionDTO.publicKey.Y)}
}
func (transactionDTO *TransactionDTO) SignatureString() [2]string {
	return transactionDTO.signature.String()
}
func (transactionDTO *TransactionDTO) GetSignature() signature.Signature {
	return *transactionDTO.signature
}
func (transactionDTO *TransactionDTO) GetPublicKey() ecdsa.PublicKey {
	return *transactionDTO.publicKey
}
func (transactionDTO *TransactionDTO) Print() {
	fmt.Printf("$	sender:		%s\n", transactionDTO.SenderAddress())
	fmt.Printf("$	recipient:	%s\n", transactionDTO.RecipientAddress())
	fmt.Printf("$	value:		%3f\n", transactionDTO.Value())
	fmt.Printf("$	signature:		%3s\n", transactionDTO.Signature)
	fmt.Printf("$	publicKey:		%3s\n", transactionDTO.PublicKey)
}

func (transactionDTO *TransactionDTO) DecodePublicKey(publicKeyString [2]string) (*ecdsa.PublicKey, error) {
	var x, y big.Int

	byteX, err := hex.DecodeString(publicKeyString[0])
	if err != nil {
		return nil, err
	}
	byteY, err := hex.DecodeString(publicKeyString[1])
	if err != nil {
		return nil, err
	}
	return &ecdsa.PublicKey{elliptic.P384(), x.SetBytes(byteX), y.SetBytes(byteY)}, nil
}
