package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

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
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, d)
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
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}
