package transaction

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float64
}


func New(senderAddress string, recipientAddress string, value float64) *Transaction {
	return &Transaction{
		senderAddress:    senderAddress,
		recipientAddress: recipientAddress,
		value:            value,
	}
}

func (transaction *Transaction) SenderAddress() string {
	return transaction.senderAddress
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) Value() float64 {
	return transaction.value
}

func (transaction *Transaction) Print() {
	fmt.Printf("$	sender:		%s\n", transaction.senderAddress)
	fmt.Printf("$	recipient:	%s\n", transaction.recipientAddress)
	fmt.Printf("$	value:		%3f\n", transaction.value)
}
func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float64 `json:"value"`
	}{
		Sender:    transaction.senderAddress,
		Recipient: transaction.recipientAddress,
		Value:     transaction.value,
	})
}
