package signature

import (
	"fmt"
	"math/big"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func New(r, s *big.Int) *Signature {
	return &Signature{r: r, s: s}
}
func (signature *Signature) GetR() big.Int {
	return *signature.r
}
func (signature *Signature) GetS() big.Int {
	return *signature.s
}
func (signature *Signature) String() string {
	return fmt.Sprintf("%x%x", signature.r, signature.s)
}
