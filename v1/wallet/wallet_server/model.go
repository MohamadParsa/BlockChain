package wallet_server

type SendCryptoRequest struct {
	RecipientAddress string `json:"recipientAddress"`
	Amount           string `json:"amount"`
}
