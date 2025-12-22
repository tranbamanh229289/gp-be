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
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
)

type IClaimService interface {
	GenerateSignatureBasedInputs(
		ctx context.Context,
		requestID *big.Int,
		holderID *core.ID,
		verifierID *core.ID,
		issuerState *IdentityState,
		claim *core.Claim,
		claimSignature *babyjub.Signature,
		query circuits.Query,
	) ([]byte, error)

	GenerateMTPBasedInputs(
		ctx context.Context,
		requestID *big.Int,
		holderID *core.ID,
		verifierID *core.ID,
		issuerState *IdentityState,
		claim *core.Claim,
		query circuits.Query,
	) ([]byte, error)

	GenerateAuthV3Inputs(
		ctx context.Context,
		userState *IdentityState,
		gistTree *merkletree.MerkleTree,
		challenge *big.Int,
		challengeSignature *babyjub.Signature,
	) ([]byte, error)

	GenerateStateTransitionInputs(
		ctx context.Context,
		oldState *IdentityState,
		newState *IdentityState,
		isOldStateGenesis bool,
		authClaimSignature *babyjub.Signature,
	) ([]byte, error)
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
	requestID *big.Int,
	holderID *core.ID,
	verifierID *core.ID,
	issuerState *IdentityState,
	claim *core.Claim,
	claimSignature *babyjub.Signature,
	query circuits.Query,
) ([]byte, error) {
	s.logger.Info("🔐 Generating Signature-Based Proof Inputs")

	// 1. calculate current state
	currentState, err := issuerState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 2. calculate tree issuerState
	treeState := circuits.TreeState{
		State:          currentState,
		ClaimsRoot:     issuerState.ClaimsTree.Root(),
		RevocationRoot: issuerState.RevocationsTree.Root(),
		RootOfRoots:    issuerState.RootsTree.Root(),
	}

	// 4. get auth claim
	authClaim, err := issuerState.GetAuthClaim()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get auth claim proof
	issuerAuthClaimProof, err := issuerState.GetClaimsTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. get non-rev auth claim proof
	issuerAuthClaimNonRevProof, err := issuerState.GetNonRevTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 8. get non-rev claim proof
	issuerClaimNonRevProof, err := issuerState.GetNonRevTreeProof(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 9. get issuerID
	did, err := issuerState.GetDID()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	issuerID, err := core.IDFromDID(did)
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
	requestID *big.Int,
	holderID *core.ID,
	verifierID *core.ID,
	issuerState *IdentityState,
	claim *core.Claim,
	query circuits.Query,
) ([]byte, error) {
	s.logger.Info("🌳 Generating MTP-Based Proof Inputs")
	// 1. calculate current state
	currentState, err := issuerState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 2. calculate tree issuerState
	treeState := circuits.TreeState{
		State:          currentState,
		ClaimsRoot:     issuerState.ClaimsTree.Root(),
		RevocationRoot: issuerState.RevocationsTree.Root(),
		RootOfRoots:    issuerState.RootsTree.Root(),
	}

	// 3. get claim proof
	issuerClaimProof, err := issuerState.GetClaimsTreeProof(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get non-rev claim proof
	issuerClaimNonRevProof, err := issuerState.GetNonRevTreeProof(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get issuer ID
	did, err := issuerState.GetDID()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	issuerID, err := core.IDFromDID(did)

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
	userState *IdentityState,
	gistTree *merkletree.MerkleTree,
	challenge *big.Int,
	challengeSignature *babyjub.Signature,
) ([]byte, error) {
	s.logger.Info("\n🔐 Generating AuthV3 Circuit Inputs")

	// 1. calculate state
	userStateValue, err := userState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 2. get auth claim
	authClaim, err := userState.GetAuthClaim()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 3. get auth claim proof
	authClaimProof, err := userState.GetClaimsTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get non-rev auth claim proof
	authClaimNonRevProof, err := userState.GetNonRevTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get state gist Proof
	gistProof, _, err := gistTree.GenerateProof(ctx, userStateValue.BigInt(), gistTree.Root())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 7. user id
	did, err := userState.GetDID()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	userID, err := core.IDFromDID(did)
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
			Proof: gistProof,
		},
		Signature: challengeSignature,
		Challenge: challenge,
	}

	return inputs.InputsMarshal()
}

func (s *ClaimService) GenerateStateTransitionInputs(
	ctx context.Context,
	oldState *IdentityState,
	newState *IdentityState,
	isOldStateGenesis bool,
	authClaimSignature *babyjub.Signature,
) ([]byte, error) {
	s.logger.Info("\n🔄 Generating StateTransition Circuit Inputs")

	// 1. calculate old state
	oldStateValue, err := oldState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 2. calculate new state
	newStateValue, err := newState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 3. get auth claim
	authClaim, err := oldState.GetAuthClaim()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get auth claim proof with old state
	authClaimOldStateProof, err := oldState.GetClaimsTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get non-rev auth claim proof with old state
	authClaimNonRevOldStateProof, err := oldState.GetNonRevTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. get auth claim proof with new state
	authClaimNewStateProof, err := oldState.GetClaimsTreeProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 8. get user id
	did, err := oldState.GetDID()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	userID, err := core.IDFromDID(did)
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
		Signature:               authClaimSignature,
	}

	return inputs.InputsMarshal()
}
