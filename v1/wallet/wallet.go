package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func New() (*Wallet, error) {
	wallet := &Wallet{}
	var err error
	wallet.privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to generate wallet")
	}
	wallet.publicKey = &wallet.privateKey.PublicKey
	return wallet, nil
}

func (wallet *Wallet) generateWalletAddress() (string, error) {
	address := ""
	// 0 - Having a private ECDSA key
	//|-> wallet.privateKey

	// 1 - Take the corresponding public key generated with it (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)
	// |-> wallet.publicKey

	// 2 - Perform SHA-256 hashing on the public key
	sha := sha256.New()
	sha.Write(wallet.publicKey.X.Bytes())
	sha.Write(wallet.publicKey.Y.Bytes())
	shaPublicKey := sha.Sum(nil)

	// 3 - Perform RIPEMD-160 hashing on the result of SHA-256
	rip := ripemd160.New()
	rip.Write(shaPublicKey)
	ripPublicKey := rip.Sum(nil)

	// 4 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	ripVersionBytes := make([]byte, 21)
	ripVersionBytes = append([]byte{0x00}, ripPublicKey...)

	// 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
	sha = sha256.New()
	sha.Write(ripVersionBytes)
	shaRipPublicKey := sha.Sum(nil)

	// 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
	sha = sha256.New()
	sha.Write(shaRipPublicKey)
	shaRipPublicKey = sha.Sum(nil)

	// 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
	checkSum := shaRipPublicKey[:4]

	// 8 - Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	ripVersionBytes = append(ripVersionBytes, checkSum...)
	// 9 - Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	address = base58.Encode(ripVersionBytes)
	return address, nil
}
func (wallet *Wallet) PublicKey() string {
	if wallet == nil {
		return ""
	}
	return fmt.Sprintf("%x%x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}

func (wallet *Wallet) Sign(text string) (string, error) {
	if wallet == nil {
		return "", errors.New("Wallet not initialized")
	}
	r, s, err := wallet.sign(text)
	if err != nil {
		return "", errors.New("failed to sign")
	}
	return fmt.Sprintf("%x%x", r, s), nil
}

func (wallet *Wallet) sign(text string) (r *big.Int, s *big.Int, err error) {
	hash := sha256.Sum256([]byte(text))
	r, s, err = ecdsa.Sign(rand.Reader, wallet.privateKey, hash[:])
	if err != nil {
		return r, s, errors.New("failed to sign")
	}
	return r, s, nil
}
func (wallet *Wallet) Verify(hash []byte) bool {
	if wallet == nil {
		return false
	}
	r, s, err := wallet
	return ecdsa.Verify(wallet.publicKey, hash)
}
