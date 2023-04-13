package signature

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
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
func (signature *Signature) String() [2]string {
	return [2]string{fmt.Sprintf("%x", signature.r), fmt.Sprintf("%x", signature.s)}
}
func (signature *Signature) DecodeSignature(signatureString [2]string) (*Signature, error) {
	var r, s big.Int
	byteR, err := hex.DecodeString(signatureString[0])
	if err != nil {
		return nil, err
	}
	byteS, err := hex.DecodeString(signatureString[1])
	if err != nil {
		return nil, err
	}
	return &Signature{r: r.SetBytes(byteR), s: s.SetBytes(byteS)}, nil
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
