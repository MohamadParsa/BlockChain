package wallet_server

type SendCryptoRequestDTO struct {
	RecipientAddress string `json:"recipientAddress"`
	Amount           string `json:"amount"`
}
