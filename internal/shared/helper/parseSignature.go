package helper

import (
	"encoding/hex"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

func GetSignatureFromString(signatureString string) (*babyjub.Signature, error) {
	signCompBytes, err := hex.DecodeString(signatureString)
	if err != nil {
		return nil, err
	}
	var compSig babyjub.SignatureComp
	copy(compSig[:], signCompBytes)
	signature, err := compSig.Decompress()
	if err != nil {
		return nil, err
	}
	return signature, nil
}
