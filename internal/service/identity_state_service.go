package service

import (
	"context"
	"fmt"
	"math/big"

	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
)

type IdentityState struct {
	PublicKey       *babyjub.PublicKey
	ClaimsTree      *merkletree.MerkleTree
	ClaimsMTID      uint64
	RevocationsTree *merkletree.MerkleTree
	RevMTID         uint64
	RootsTree       *merkletree.MerkleTree
	RootsMTID       uint64
}

func (state *IdentityState) GetStateValue() (*merkletree.Hash, error) {
	stateValue, err := merkletree.HashElems(state.ClaimsTree.Root().BigInt(), state.RevocationsTree.Root().BigInt(), state.RootsTree.Root().BigInt())

	return stateValue, err
}

func (state *IdentityState) GetDID() (w3c.DID, error) {
	stateValue, _ := state.GetStateValue()
	typ, _ := core.BuildDIDType(core.DIDMethodPolygonID, core.Polygon, core.Mumbai)
	did, _ := core.NewDIDFromIdenState(typ, stateValue.BigInt())
	return *did, nil
}

func (state *IdentityState) AddClaim(ctx context.Context, claim *core.Claim) error {
	hi, hv, err := claim.HiHv()
	if err != nil {
		return fmt.Errorf("failed to get key and value: %w", err)
	}
	err = state.ClaimsTree.Add(ctx, hi, hv)
	if err != nil {
		return fmt.Errorf("failed to add claim: %w", err)
	}

	return nil
}

func (state *IdentityState) RevokeClaim(ctx context.Context, claim *core.Claim) error {
	revNonce := claim.GetRevocationNonce()
	err := state.RevocationsTree.Add(ctx, new(big.Int).SetUint64(revNonce), big.NewInt(0))
	if err != nil {
		return fmt.Errorf("failed to add claim: %w", err)
	}

	return nil
}

func (state *IdentityState) GetAuthClaim() (*core.Claim, error) {
	authSchemaHash, _ := core.NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	revNonce := uint64(1)

	authClaim, err := core.NewClaim(
		authSchemaHash,
		core.WithIndexDataInts(state.PublicKey.X, state.PublicKey.Y),
		core.WithRevocationNonce(revNonce),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create auth claim: %w", err)
	}

	return authClaim, nil
}

func (state *IdentityState) GetClaimsTreeProof(ctx context.Context, claim *core.Claim) (*merkletree.Proof, error) {
	hi, _, err := claim.HiHv()
	if err != nil {
		return nil, fmt.Errorf("failed to get key and value: %w", err)
	}

	proof, _, err := state.ClaimsTree.GenerateProof(ctx, hi, state.ClaimsTree.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth claim proof: %w", err)
	}

	return proof, nil
}

func (state *IdentityState) GetNonRevTreeProof(ctx context.Context, claim *core.Claim) (*merkletree.Proof, error) {
	revNonce := claim.GetRevocationNonce()
	proof, _, err := state.RevocationsTree.GenerateProof(ctx, new(big.Int).SetUint64(revNonce), state.RevocationsTree.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to generate non-rev proof: %w", err)
	}

	return proof, nil
}

// func (state *IdentityState) SignClaim(claim *core.Claim) *babyjub.Signature {
// 	index, value := claim.RawSlots()
// 	indexHash, _ := poseidon.Hash(core.ElemBytesToInts(index[:]))
// 	valueHash, _ := poseidon.Hash(core.ElemBytesToInts(value[:]))
// 	claimHash, _ := merkletree.HashElems(indexHash, valueHash)

// 	return state.privateKey.SignPoseidon(claimHash.BigInt())
// }

// func (state *IdentityState) SignChallenge(challenge *big.Int) *babyjub.Signature {
// 	return state.privateKey.SignPoseidon(challenge)
// }
