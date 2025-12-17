package utils

import (
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-iden3-crypto/keccak256"
)

func GetSchemaHash(url string) (*core.SchemaHash, error) {
	var sHash core.SchemaHash
	h := keccak256.Hash([]byte(url))
	copy(sHash[:], h[len(h)-16:])
	sHashIndex, err := sHash.MarshalText()
	if err != nil {
		return nil, err
	}
	claim, _ := core.NewSchemaHashFromHex(string(sHashIndex))
	return &claim, nil
}

func GetIDFromIDI(DID *w3c.DID) string {
	return DID.ID
}
