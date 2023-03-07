package transaction

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	SenderAddress    string  `json:"senderAddress"`
	RecipientAddress string  `json:"recipientAddress"`
	Value            float64 `json:"value"`
}

func New(senderAddress string, recipientAddress string, value float64) *Transaction {
	return &Transaction{
		SenderAddress:    senderAddress,
		RecipientAddress: recipientAddress,
		Value:            value,
	}
}

func (transaction *Transaction) Print() {
	fmt.Printf("$	sender:		%s\n", transaction.SenderAddress)
	fmt.Printf("$	recipient:	%s\n", transaction.RecipientAddress)
	fmt.Printf("$	value:		%3f\n", transaction.Value)
}
func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"senderAddress"`
		Recipient string  `json:"recipientAddress"`
		Value     float64 `json:"value"`
	}{
		Sender:    transaction.SenderAddress,
		Recipient: transaction.RecipientAddress,
		Value:     transaction.Value,
	})
}
