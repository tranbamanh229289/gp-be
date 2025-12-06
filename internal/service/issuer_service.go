package service

import (
	"context"
	"math/big"
	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
)

type IIssuerService interface {
	CreateIssuerIdentity(ctx context.Context, privateKey *babyjub.PrivateKey) (*core.DID, *core.ID)
}

type IssuerService struct {
}

type IdentityState struct {
	ClaimsTree      *merkletree.MerkleTree
	RevocationsTree *merkletree.MerkleTree
	RootsTree       *merkletree.MerkleTree
	State           *merkletree.Hash
	PublicKey       *babyjub.PublicKey
	DID             *core.DID
}

func (is *IssuerService) CreateIssuerIdentity(ctx context.Context, privateKey *babyjub.PrivateKey) *IdentityState {
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

func (is *IssuerService) SignClaim(privKey *babyjub.PrivateKey, claim *core.Claim) *babyjub.Signature {
	index, value := claim.RawSlots()
	indexHash, _ := poseidon.Hash(core.ElemBytesToInts(index[:]))
	valueHash, _ := poseidon.Hash(core.ElemBytesToInts(value[:]))
	claimHash, _ := merkletree.HashElems(indexHash, valueHash)

	return privKey.SignPoseidon(claimHash.BigInt())
}

func (is *IdentityState) CreateClaim(ctx context.Context, schemaHash string, subjectID core.ID, indexSlotA,
	indexSlotB *big.Int, revNonce uint64, expireDate time.Time) error {
	// update root tree
	is.RootsTree.Add(ctx, is.ClaimsTree.Root().BigInt(), big.NewInt(0))

	// create claim
	claimHash, _ := core.NewSchemaHashFromHex(schemaHash)
	claim, _ := core.NewClaim(claimHash, core.WithIndexID(subjectID),
		core.WithIndexDataInts(indexSlotA, indexSlotB),
		core.WithRevocationNonce(revNonce),
		core.WithExpirationDate(expireDate))
	claimHi, claimHv, _ := claim.HiHv()

	// add claim
	is.ClaimsTree.Add(ctx, claimHi, claimHv)

	// update new state
	newState, _ := merkletree.HashElems(is.ClaimsTree.Root().BigInt(), is.RevocationsTree.Root().BigInt(), is.RootsTree.Root().BigInt())
	is.State = newState
	return nil
}

func (is *IdentityState) RevokeClaim(ctx context.Context, revNonce uint64) error {
	// add claim
	is.RevocationsTree.Add(ctx, new(big.Int).SetUint64(revNonce), big.NewInt(0))

	// update new state
	newState, _ := merkletree.HashElems(is.ClaimsTree.Root().BigInt(), is.RevocationsTree.Root().BigInt(), is.RootsTree.Root().BigInt())
	is.State = newState
	return nil
}
