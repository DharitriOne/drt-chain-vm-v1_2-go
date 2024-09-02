package bls

import (
	crypto "github.com/DharitriOne/drt-chain-crypto-go"
	"github.com/DharitriOne/drt-chain-crypto-go/signing"
	"github.com/DharitriOne/drt-chain-crypto-go/signing/mcl"
	"github.com/DharitriOne/drt-chain-crypto-go/signing/mcl/singlesig"
)

type bls struct {
	keyGenerator crypto.KeyGenerator
	signer       crypto.SingleSigner
}

func NewBLS() *bls {
	b := &bls{}
	suite := mcl.NewSuiteBLS12()
	b.keyGenerator = signing.NewKeyGenerator(suite)
	b.signer = singlesig.NewBlsSigner()

	return b
}

func (b *bls) VerifyBLS(key []byte, msg []byte, sig []byte) error {
	publicKey, err := b.keyGenerator.PublicKeyFromByteArray(key)
	if err != nil {
		return err
	}

	return b.signer.Verify(publicKey, msg, sig)
}
