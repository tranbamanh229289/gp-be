package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/iden3/go-merkletree-sql/v2"
)

type Hash string

func (h Hash) ToMerkleTreeHash() (*merkletree.Hash, error) {
	var mth merkletree.Hash
	err := mth.UnmarshalText([]byte(h))
	return &mth, err
}

type DID string

func (d *DID) Scan(value interface{}) error {
	if value == nil {
		*d = ""
		return nil
	}
	return json.Unmarshal(value.([]byte), d)
}

func (d DID) Value() (driver.Value, error) {
	return string(d), nil
}

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal(value.([]byte), j)
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}
