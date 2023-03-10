package signature

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/MohamadParsa/BlockChain/v1/transaction"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func New(r, s *big.Int) *Signature {
	return &Signature{r: r, s: s}
}

func (signature *Signature) GetR() *big.Int {
	r := *signature.r
	return &r
}
func (signature *Signature) GetS() *big.Int {
	s := *signature.s
	return &s
}
func (signature *Signature) String() string {
	return fmt.Sprintf("%x%x", signature.r, signature.s)
}
func VerifySignature(publicKey *ecdsa.PublicKey, signature *Signature, transaction *transaction.Transaction) (bool, error) {
	if signature == nil {
		return false, errors.New("signature is invalid")
	}
	transactionJsonBytes, err := json.Marshal(transaction)
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256(transactionJsonBytes)
	return ecdsa.Verify(publicKey, hash[:], signature.GetR(), signature.GetS()), nil
}
