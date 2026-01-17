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
	PublicKey  *babyjub.PublicKey
	DID        *w3c.DID
	ClaimsTree *merkletree.MerkleTree
	ClaimsMTID uint64
	RevTree    *merkletree.MerkleTree
	RevMTID    uint64
	RootsTree  *merkletree.MerkleTree
	RootsMTID  uint64
}

func (state *IdentityState) GetStateValue() (*merkletree.Hash, error) {
	return merkletree.HashElems(state.ClaimsTree.Root().BigInt(), state.RevTree.Root().BigInt(), state.RootsTree.Root().BigInt())

}

func (state *IdentityState) GetDID() *w3c.DID {
	return state.DID
}

func (state *IdentityState) GetID() *core.ID {
	did, _ := core.IDFromDID(*state.DID)
	return &did
}

func (state *IdentityState) CalculateState(claimsTree, revTree, rootsTree *merkletree.MerkleTree) (*merkletree.Hash, error) {
	return merkletree.HashElems(claimsTree.Root().BigInt(), revTree.Root().BigInt(), rootsTree.Root().BigInt())
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

	claimsRoot := state.ClaimsTree.Root()
	err = state.RootsTree.Add(ctx, claimsRoot.BigInt(), big.NewInt(1))
	return nil
}

func (state *IdentityState) RevokeClaim(ctx context.Context, claim *core.Claim) error {
	revNonce := claim.GetRevocationNonce()
	err := state.RevTree.Add(ctx, new(big.Int).SetUint64(revNonce), big.NewInt(0))
	if err != nil {
		return fmt.Errorf("failed to add claim: %w", err)
	}
	return nil
}

func (state *IdentityState) GetAuthClaim() (*core.Claim, error) {
	revNonce := uint64(1)

	authClaim, err := core.NewClaim(
		core.AuthSchemaHash,
		core.WithIndexDataInts(state.PublicKey.X, state.PublicKey.Y),
		core.WithRevocationNonce(revNonce),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create auth claim: %w", err)
	}

	return authClaim, nil
}

func (state *IdentityState) GetIncMTProof(ctx context.Context, claim *core.Claim) (*merkletree.Proof, error) {
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

func (state *IdentityState) GetNonRevMTProof(ctx context.Context, claim *core.Claim) (*merkletree.Proof, error) {
	revNonce := claim.GetRevocationNonce()
	proof, _, err := state.RevTree.GenerateProof(ctx, new(big.Int).SetUint64(revNonce), state.RevTree.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to generate non-rev proof: %w", err)
	}

	return proof, nil
}
