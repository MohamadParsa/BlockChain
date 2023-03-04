package wallet_server

type SendCryptoRequest struct {
	RecipientAddress string `json:"recipient_address"`
	Amount           string `json:"amount"`
}
