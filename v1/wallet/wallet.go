package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"

	"github.com/MohamadParsa/BlockChain/v1/blocks/transaction"
	"github.com/MohamadParsa/BlockChain/v1/wallet/signature"
	"github.com/btcsuite/btcd/btcutil/base58"
	logger "go.uber.org/zap"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
	address    string
}

func New() (*Wallet, error) {
	wallet := &Wallet{}
	var err error
	wallet.privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to create wallet")
	}
	wallet.publicKey = &wallet.privateKey.PublicKey
	wallet.address, err = wallet.generateWalletAddress()
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to create wallet")
	}
	return wallet, nil
}

func (wallet *Wallet) generateWalletAddress() (string, error) {
	address := ""
	var lastError error
	// 0 - Having a private ECDSA key
	//|-> wallet.privateKey

	// 1 - Take the corresponding public key generated with it (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)
	// |-> wallet.publicKey

	// 2 - Perform SHA-256 hashing on the public key
	sha := sha256.New()
	_, err := sha.Write(wallet.publicKey.X.Bytes())
	logError(err, &lastError)
	_, err = sha.Write(wallet.publicKey.Y.Bytes())
	shaPublicKey := sha.Sum(nil)
	logError(err, &lastError)

	// 3 - Perform RIPEMD-160 hashing on the result of SHA-256
	rip := ripemd160.New()
	_, err = rip.Write(shaPublicKey)
	logError(err, &lastError)
	ripPublicKey := rip.Sum(nil)

	// 4 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	ripVersionBytes := make([]byte, 21)
	ripVersionBytes = append([]byte{0x00}, ripPublicKey...)

	// 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
	sha = sha256.New()
	_, err = sha.Write(ripVersionBytes)
	logError(err, &lastError)
	shaRipPublicKey := sha.Sum(nil)

	// 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
	sha = sha256.New()
	_, err = sha.Write(shaRipPublicKey)
	logError(err, &lastError)
	shaRipPublicKey = sha.Sum(nil)

	// 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
	checkSum := shaRipPublicKey[:4]

	// 8 - Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	ripVersionBytes = append(ripVersionBytes, checkSum...)
	// 9 - Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	address = base58.Encode(ripVersionBytes)
	if lastError != nil {
		return "", lastError
	}
	return address, nil
}

func (wallet *Wallet) PublicKey() string {
	if wallet == nil {
		return ""
	}
	return fmt.Sprintf("%x%x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}

func (wallet *Wallet) Address() string {
	if wallet == nil {
		return ""
	}
	return wallet.address
}
func (wallet *Wallet) SendCrypto(recipient_address string, value float64) (*transaction.Transaction, error) {
	if wallet == nil {
		return nil, errors.New("wallet is invalid")
	}
	transaction := transaction.New(wallet.Address(), recipient_address, value)
	return transaction, nil
}

func (wallet *Wallet) Sign(transaction *transaction.Transaction) (*signature.Signature, error) {
	if wallet == nil {
		return nil, errors.New("Wallet not initialized")
	}
	if transaction.SenderAddress() != wallet.Address() {
		return nil, errors.New("transaction is invalid")
	}
	transactionJsonBytes, err := transaction.MarshalJSON()
	if err != nil {
		return nil, errors.New("failed to sign")
	}
	signature, err := wallet.sign(transactionJsonBytes)
	if err != nil {
		return nil, errors.New("failed to sign")
	}
	return signature, nil
}

func (wallet *Wallet) sign(transaction []byte) (*signature.Signature, error) {
	hash := sha256.Sum256(transaction)
	r, s, err := ecdsa.Sign(rand.Reader, wallet.privateKey, hash[:])
	if err != nil {
		return nil, errors.New("failed to sign")
	}

	return signature.New(r, s), nil
}

func logError(err error, lastError *error) {
	if err != nil {
		logger.Error(err)
		*lastError = err
	}
}
