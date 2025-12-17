package service

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/iden3/go-circuits/v2"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
)

type IdentityState struct {
	DID             *w3c.DID
	PublicKey       *babyjub.PublicKey
	ClaimsTree      *merkletree.MerkleTree
	RevocationsTree *merkletree.MerkleTree
	RootsTree       *merkletree.MerkleTree
	State           *merkletree.Hash
}

type ClaimService struct {
	config *config.Config
	logger *logger.ZapLogger
}

func NewClaimService(config *config.Config, logger *logger.ZapLogger) *ClaimService {
	return &ClaimService{
		config: config,
		logger: logger,
	}
}

func (s *ClaimService) GenerateSignatureBasedInputs(
	ctx context.Context,
	privateKey *babyjub.PrivateKey,
	issuerState *IdentityState,
	claim *core.Claim,
	query circuits.Query,
	requestID *big.Int,
	holderID *core.ID,
	verifierID *core.ID,
) ([]byte, error) {
	s.logger.Info("🔐 Generating Signature-Based Proof Inputs")

	// 1. calculate current state
	currentState := issuerState.State

	// 2. calculate tree issuerState
	treeState := circuits.TreeState{
		State:          currentState,
		ClaimsRoot:     issuerState.ClaimsTree.Root(),
		RevocationRoot: issuerState.RevocationsTree.Root(),
		RootOfRoots:    issuerState.RootsTree.Root(),
	}

	// 3. sign claim signature
	claimSignature := s.signClaim(privateKey, claim)

	// 4. get auth claim
	authClaim, err := s.getAuthClaim(ctx, privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get auth claim proof
	issuerAuthClaimProof, err := s.getClaimsTreeProof(ctx, authClaim, issuerState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. get non-rev auth claim proof
	issuerAuthClaimNonRevProof, err := s.getNonRevTreeProof(ctx, authClaim, issuerState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 8. get non-rev claim proof
	issuerClaimNonRevProof, err := s.getNonRevTreeProof(ctx, claim, issuerState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 9. get issuerID
	issuerID, err := core.IDFromDID(*issuerState.DID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	inputs := circuits.AtomicQueryV3Inputs{
		RequestID:                requestID,
		ID:                       holderID,
		ProfileNonce:             big.NewInt(0),
		ClaimSubjectProfileNonce: big.NewInt(0),

		Claim: circuits.ClaimWithSigAndMTPProof{
			IssuerID: &issuerID,
			Claim:    claim,
			NonRevProof: circuits.MTProof{
				Proof:     issuerClaimNonRevProof,
				TreeState: treeState,
			},
			IncProof: &circuits.MTProof{
				Proof:     &merkletree.Proof{},
				TreeState: treeState,
			},
			SignatureProof: &circuits.BJJSignatureProof{
				Signature:       claimSignature,
				IssuerAuthClaim: authClaim,
				IssuerAuthIncProof: circuits.MTProof{
					Proof:     issuerAuthClaimProof,
					TreeState: treeState,
				},
				IssuerAuthNonRevProof: circuits.MTProof{
					Proof:     issuerAuthClaimNonRevProof,
					TreeState: treeState,
				},
			},
		},

		SkipClaimRevocationCheck: false,
		Query:                    query,
		CurrentTimeStamp:         time.Now().Unix(),
		ProofType:                circuits.BJJSignatureProofType,
		VerifierID:               verifierID,
		NullifierSessionID:       big.NewInt(0),
		LinkNonce:                big.NewInt(0),
	}
	s.logger.Info("✅ Signature-based inputs generated successfully")
	return inputs.InputsMarshal()
}

func (s *ClaimService) GenerateMTPBasedInputs(
	ctx context.Context,
	issuerState *IdentityState,
	claim *core.Claim,
	query circuits.Query,
	requestID *big.Int,
	holderID *core.ID,
	verifierID *core.ID,
) ([]byte, error) {
	s.logger.Info("🌳 Generating MTP-Based Proof Inputs")
	// 1. calculate current state
	currentState := issuerState.State

	// 2. calculate tree issuerState
	treeState := circuits.TreeState{
		State:          currentState,
		ClaimsRoot:     issuerState.ClaimsTree.Root(),
		RevocationRoot: issuerState.RevocationsTree.Root(),
		RootOfRoots:    issuerState.RootsTree.Root(),
	}

	// 3. get claim proof
	issuerClaimProof, err := s.getClaimsTreeProof(ctx, claim, issuerState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get non-rev claim proof
	issuerClaimNonRevProof, err := s.getNonRevTreeProof(ctx, claim, issuerState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get issuer ID
	issuerID, err := core.IDFromDID(*issuerState.DID)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	inputs := circuits.AtomicQueryV3Inputs{
		RequestID:                requestID,
		ID:                       holderID,
		ProfileNonce:             big.NewInt(0),
		ClaimSubjectProfileNonce: big.NewInt(0),

		Claim: circuits.ClaimWithSigAndMTPProof{
			IssuerID: &issuerID,
			Claim:    claim,
			NonRevProof: circuits.MTProof{
				Proof:     issuerClaimNonRevProof,
				TreeState: treeState,
			},
			IncProof: &circuits.MTProof{
				Proof:     issuerClaimProof,
				TreeState: treeState,
			},
			SignatureProof: nil,
		},

		SkipClaimRevocationCheck: false,
		Query:                    query,
		CurrentTimeStamp:         time.Now().Unix(),
		ProofType:                circuits.Iden3SparseMerkleTreeProofType,
		VerifierID:               verifierID,
		NullifierSessionID:       big.NewInt(0),
		LinkNonce:                big.NewInt(0),
	}

	s.logger.Info("✅ MTP-based inputs generated successfully")
	return inputs.InputsMarshal()
}

func (s *ClaimService) GenerateAuthV3Inputs(
	ctx context.Context,
	privateKey *babyjub.PrivateKey,
	userState *IdentityState,
	challenge *big.Int,
	gistTree *merkletree.MerkleTree,
) ([]byte, error) {
	s.logger.Info("\n🔐 Generating AuthV3 Circuit Inputs")

	// 1. calculate state
	userStateValue := userState.State

	// 2. get auth claim
	authClaim, err := s.getAuthClaim(ctx, privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 3. get auth claim proof
	authClaimProof, err := s.getClaimsTreeProof(ctx, authClaim, userState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get non-rev auth claim proof
	authClaimNonRevProof, err := s.getNonRevTreeProof(ctx, authClaim, userState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get state gist Proof
	authClaimGistProof, err := s.getGistTreeProof(ctx, userStateValue, gistTree)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. signature
	signature := s.signChallenge(privateKey, challenge)

	// 7. user id
	userID, err := core.IDFromDID(*userState.DID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	inputs := circuits.AuthV3Inputs{
		GenesisID:          &userID,
		ProfileNonce:       big.NewInt(0),
		AuthClaim:          authClaim,
		AuthClaimIncMtp:    authClaimProof,
		AuthClaimNonRevMtp: authClaimNonRevProof,
		TreeState: circuits.TreeState{
			State:          userStateValue,
			ClaimsRoot:     userState.ClaimsTree.Root(),
			RevocationRoot: userState.RevocationsTree.Root(),
			RootOfRoots:    userState.RootsTree.Root(),
		},
		GISTProof: circuits.GISTProof{
			Root:  gistTree.Root(),
			Proof: authClaimGistProof,
		},
		Signature: signature,
		Challenge: challenge,
	}

	return inputs.InputsMarshal()
}

func (s *ClaimService) GenerateStateTransitionInputs(
	ctx context.Context,
	privateKey *babyjub.PrivateKey,
	oldState *IdentityState,
	newState *IdentityState,
	isOldStateGenesis bool,
) ([]byte, error) {
	s.logger.Info("\n🔄 Generating StateTransition Circuit Inputs")

	// 1. calculate old state
	oldStateValue := oldState.State

	// 2. calculate new state
	newStateValue := newState.State

	// 3. get auth claim
	authClaim, err := s.getAuthClaim(ctx, privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get auth claim proof with old state
	authClaimOldStateProof, err := s.getClaimsTreeProof(ctx, authClaim, oldState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get non-rev auth claim proof with old state
	authClaimNonRevOldStateProof, err := s.getNonRevTreeProof(ctx, authClaim, oldState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. get auth claim proof with new state
	authClaimNewStateProof, err := s.getClaimsTreeProof(ctx, authClaim, newState)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 7. get auth claim proof with new state
	signature := s.signClaim(privateKey, authClaim)

	// 8. get user id
	userID, err := core.IDFromDID(*oldState.DID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	inputs := circuits.StateTransitionInputs{
		ID: &userID,
		OldTreeState: circuits.TreeState{
			State:          oldStateValue,
			ClaimsRoot:     oldState.ClaimsTree.Root(),
			RevocationRoot: oldState.RevocationsTree.Root(),
			RootOfRoots:    oldState.RootsTree.Root(),
		},
		NewTreeState: circuits.TreeState{
			State:          newStateValue,
			ClaimsRoot:     newState.ClaimsTree.Root(),
			RevocationRoot: newState.RevocationsTree.Root(),
			RootOfRoots:    newState.RootsTree.Root(),
		},
		IsOldStateGenesis:       isOldStateGenesis,
		AuthClaim:               authClaim,
		AuthClaimIncMtp:         authClaimOldStateProof,
		AuthClaimNonRevMtp:      authClaimNonRevOldStateProof,
		AuthClaimNewStateIncMtp: authClaimNewStateProof,
		Signature:               signature,
	}

	return inputs.InputsMarshal()
}

func (s *ClaimService) CreateGenesisState(ctx context.Context, publicKey *babyjub.PublicKey) *IdentityState {

	// create auth claim
	authSchemaHash, _ := core.NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	revNonce := uint64(1)
	authClaim, _ := core.NewClaim(authSchemaHash, core.WithIndexDataInts(publicKey.X, publicKey.Y), core.WithRevocationNonce(revNonce))

	// add auth claim into claimsTree
	storage := 
	claimsTree, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), s.config.Circuit.MTLevel)
	authHi, authHv, _ := authClaim.HiHv()
	claimsTree.Add(ctx, authHi, authHv)

	// create empty revocationsTree
	revocationsTree, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), s.config.Circuit.MTLevel)

	// create empty rootsTree
	rootsTree, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), s.config.Circuit.MTLevel)

	state, _ := s.getState(claimsTree, revocationsTree, rootsTree)

	typ, _ := core.BuildDIDType(core.DIDMethodPolygonID, core.Polygon, core.Mumbai)
	did, _ := core.NewDIDFromIdenState(typ, state.BigInt())

	return &IdentityState{
		ClaimsTree:      claimsTree,
		RevocationsTree: revocationsTree,
		RootsTree:       rootsTree,
		State:           state,
		PublicKey:       publicKey,
		DID:             did,
	}
}

func (s *ClaimService) CreateClaim(ctx context.Context, identityState *IdentityState, schemaHash string, subjectID core.ID, indexSlotA,
	indexSlotB *big.Int, revNonce uint64, expireDate time.Time) error {
	// update root tree
	identityState.RootsTree.Add(ctx, identityState.ClaimsTree.Root().BigInt(), big.NewInt(0))

	// create claim
	claimHash, _ := core.NewSchemaHashFromHex(schemaHash)
	claim, _ := core.NewClaim(claimHash, core.WithIndexID(subjectID),
		core.WithIndexDataInts(indexSlotA, indexSlotB),
		core.WithRevocationNonce(revNonce),
		core.WithExpirationDate(expireDate))
	claimHi, claimHv, _ := claim.HiHv()

	// add claim
	identityState.ClaimsTree.Add(ctx, claimHi, claimHv)

	return nil
}

func (s *ClaimService) RevokeClaim(ctx context.Context, identityState *IdentityState, revNonce uint64) error {
	// add claim
	identityState.RevocationsTree.Add(ctx, new(big.Int).SetUint64(revNonce), big.NewInt(0))

	return nil
}

func (c *ClaimService) getAuthClaim(ctx context.Context, pubKey *babyjub.PublicKey) (*core.Claim, error) {
	authSchemaHash, _ := core.NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	revNonce := uint64(1)

	authClaim, err := core.NewClaim(
		authSchemaHash,
		core.WithIndexDataInts(pubKey.X, pubKey.Y),
		core.WithRevocationNonce(revNonce),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create auth claim: %w", err)
	}

	return authClaim, nil
}

func (s *ClaimService) signClaim(privateKey *babyjub.PrivateKey, claim *core.Claim) *babyjub.Signature {
	index, value := claim.RawSlots()
	indexHash, _ := poseidon.Hash(core.ElemBytesToInts(index[:]))
	valueHash, _ := poseidon.Hash(core.ElemBytesToInts(value[:]))
	claimHash, _ := merkletree.HashElems(indexHash, valueHash)

	return privateKey.SignPoseidon(claimHash.BigInt())
}

func (s *ClaimService) signChallenge(privateKey *babyjub.PrivateKey, challenge *big.Int) *babyjub.Signature {
	return privateKey.SignPoseidon(challenge)
}

func (c *ClaimService) getClaimsTreeProof(ctx context.Context, claim *core.Claim, state *IdentityState) (*merkletree.Proof, error) {
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

func (c *ClaimService) getNonRevTreeProof(ctx context.Context, claim *core.Claim, state *IdentityState) (*merkletree.Proof, error) {
	revNonce := claim.GetRevocationNonce()
	proof, _, err := state.RevocationsTree.GenerateProof(ctx, new(big.Int).SetUint64(revNonce), state.RevocationsTree.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to generate non-rev proof: %w", err)
	}
	return proof, nil
}

func (c *ClaimService) getGistTreeProof(ctx context.Context, state *merkletree.Hash, gistTree *merkletree.MerkleTree) (*merkletree.Proof, error) {
	proof, _, err := gistTree.GenerateProof(ctx, state.BigInt(), gistTree.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth gist proof %w", err)
	}
	return proof, nil
}

func (c *ClaimService) getState(claimsTree, revocationsTree, rootsTree *merkletree.MerkleTree) (*merkletree.Hash, error) {
	state, err := merkletree.HashElems(claimsTree.Root().BigInt(), revocationsTree.Root().BigInt(), rootsTree.Root().BigInt())
	return state, err
}
