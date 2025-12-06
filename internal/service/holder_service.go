package service

import (
	"context"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
)

type IHolderService interface {
	CreateHolderIdentity(ctx context.Context, privateKey *babyjub.PrivateKey) (*core.DID, *core.ID)
}

func (is *IssuerService) CreateHolderIdentity(ctx context.Context, privateKey *babyjub.PrivateKey) *IdentityState {
	pubKey := privateKey.Public()
	authSchemaHash, _ := core.NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	revNonce := uint64(1)
	authClaim, _ := core.NewClaim(authSchemaHash, core.WithIndexDataInts(pubKey.X, pubKey.Y), core.WithRevocationNonce(revNonce))

	claimsTree, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	authHi, authHv, _ := authClaim.HiHv()
	claimsTree.Add(ctx, authHi, authHv)

	revocationsTree, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	rootsTree, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	state, _ := merkletree.HashElems(claimsTree.Root().BigInt(), revocationsTree.Root().BigInt(), rootsTree.Root().BigInt())

	typ, _ := core.BuildDIDType(core.DIDMethodPolygonID, core.Polygon, core.Mumbai)
	did, _ := core.DIDGenesisFromIdenState(typ, state.BigInt())
	return &IdentityState{
		ClaimsTree:      claimsTree,
		RevocationsTree: revocationsTree,
		RootsTree:       rootsTree,
		State:           state,
		PublicKey:       pubKey,
		DID:             did,
	}
}
