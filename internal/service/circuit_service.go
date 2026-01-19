package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/proof"
	"be/internal/shared/constant"
	"be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/iden3/go-circuits/v2"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-schema-processor/v2/loaders"
	"github.com/iden3/go-schema-processor/v2/merklize"
	"github.com/iden3/go-schema-processor/v2/verifiable"
	"gorm.io/gorm"
)

type ICircuitService interface {
	GetCredentialAtomicQueryV3Input(ctx context.Context, request *dto.CredentialAtomicQueryV3InputRequestDto, claims *dto.ZKClaims) (map[string]interface{}, error)
}
type CircuitService struct {
	config          *config.Config
	logger          *logger.ZapLogger
	proofRepo       proof.IProofRepository
	vcRepo          credential.IVerifiableCredentialRepository
	identityService IIdentityService
}

func NewCircuitService(
	config *config.Config,
	logger *logger.ZapLogger,
	proofRepo proof.IProofRepository,
	vcRepo credential.IVerifiableCredentialRepository,
	identityService IIdentityService) ICircuitService {
	return &CircuitService{
		config:          config,
		logger:          logger,
		proofRepo:       proofRepo,
		vcRepo:          vcRepo,
		identityService: identityService,
	}
}

func (s *CircuitService) GetCredentialAtomicQueryV3Input(ctx context.Context, request *dto.CredentialAtomicQueryV3InputRequestDto, claims *dto.ZKClaims) (map[string]interface{}, error) {
	proofRequest, err := s.proofRepo.FindProofRequestByPublicId(ctx, request.ProofRequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofRequestNotFound
		}
		return nil, &constant.InternalServer
	}

	vc, err := s.vcRepo.FindVerifiableCredentialByCredentialId(ctx, request.CredentialID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofRequestNotFound
		}
		return nil, &constant.InternalServer
	}

	if claims.DID != vc.HolderDID {
		return nil, &constant.BadRequest
	}

	var credentialSubject map[string]interface{}
	credSubBytes, _ := json.Marshal(proofRequest.CredentialSubject)
	json.Unmarshal(credSubBytes, &credentialSubject)

	credentialSubject = helper.NormalizeToIntMap(credentialSubject)

	query := &pubsignals.Query{
		AllowedIssuers:           proofRequest.AllowedIssuers,
		CredentialSubject:        credentialSubject,
		Context:                  proofRequest.Schema.ContextURL,
		Type:                     proofRequest.Schema.Type,
		SkipClaimRevocationCheck: proofRequest.SkipClaimRevocationCheck,
		ProofType:                proofRequest.ProofType,
		GroupID:                  proofRequest.GroupID,
	}

	w3c := dto.ToW3CCredential(vc)
	input, err := s.generateCredentialAtomicQueryV3(ctx, w3c, query, proofRequest.VerifierDID, vc.Signature)
	if err != nil {
		fmt.Println(err)
		return nil, &constant.InternalServer
	}
	var circuitInputs map[string]interface{}
	if err := json.Unmarshal(input, &circuitInputs); err != nil {
		return nil, &constant.InternalServer
	}
	return circuitInputs, nil
}

func (s *CircuitService) generateCredentialAtomicQueryV3(
	ctx context.Context,
	vc *verifiable.W3CCredential,
	query *pubsignals.Query,
	verifierDIDString string,
	signatureString string,
) ([]byte, error) {
	// Get issuer state
	issuerState, err := s.identityService.GetIdentityStateByDID(ctx, vc.Issuer)
	if err != nil {

		return nil, err
	}

	// Generate request id
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	requestId, err := rand.Int(rand.Reader, max)
	if err != nil {

		return nil, err
	}

	// Get user ID from credential subject
	raw, ok := vc.CredentialSubject["id"]
	if !ok {

		return nil, fmt.Errorf("credentialSubject not existed id")
	}
	id, ok := raw.(string)
	if !ok {

		return nil, fmt.Errorf("id is invalid")
	}

	userDID, err := w3c.ParseDID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user DID: %w", err)
	}
	userId, err := core.IDFromDID(*userDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID from DID: %w", err)
	}

	verifierDID, err := w3c.ParseDID(verifierDIDString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse verifier DID: %w", err)
	}
	verifierId, err := core.IDFromDID(*verifierDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get verifier ID from DID: %w", err)
	}

	issuerDID, err := w3c.ParseDID(vc.Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issuer DID: %w", err)
	}
	issuerId, err := core.IDFromDID(*issuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issuer ID from DID: %w", err)
	}

	// Get core claim based on proof type
	var claim *core.Claim
	switch query.ProofType {
	case string(circuits.Iden3SparseMerkleTreeProofType):

		if len(vc.Proof) < 1 {
			return nil, fmt.Errorf("vc.Proof[0] not found for MTP proof type")
		}
		claim, err = vc.Proof[0].GetCoreClaim()
		if err != nil {
			return nil, err
		}
	case string(circuits.BJJSignatureProofType):
		if len(vc.Proof) < 2 {

			return nil, fmt.Errorf("vc.Proof[1] not found for BJJ signature proof type")
		}
		claim, err = vc.Proof[1].GetCoreClaim()
		if err != nil {

			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid proof type: %s", query.ProofType)
	}

	// Get auth claim
	authClaim, err := issuerState.GetAuthClaim()
	if err != nil {
		return nil, err
	}

	// Get claim merkle tree proofs
	claimIncMtp, err := issuerState.GetIncMTProof(ctx, claim)
	if err != nil {

		return nil, err
	}

	claimNonRevMtp, err := issuerState.GetNonRevMTProof(ctx, claim)
	if err != nil {
		return nil, err
	}

	// Get auth claim merkle tree proofs
	authClaimIncMtp, err := issuerState.GetIncMTProof(ctx, authClaim)
	if err != nil {
		return nil, err
	}

	authClaimNonRevMtp, err := issuerState.GetNonRevMTProof(ctx, authClaim)
	if err != nil {
		return nil, err
	}

	// Get issuer state value
	issuerStateValue, err := issuerState.GetStateValue()
	if err != nil {
		return nil, err
	}

	// Parse signature
	signature, err := helper.GetSignatureFromString(signatureString)
	if err != nil {
		return nil, err
	}

	documentLoader := loaders.NewDocumentLoader(nil, "")

	remoteDoc, err := documentLoader.LoadDocument(query.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to load context: %w", err)
	}

	ldContextBytes, err := json.Marshal(remoteDoc.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	merklizerOptions := merklize.Options{
		DocumentLoader: documentLoader,
	}

	ldContextJSON := string(ldContextBytes)
	metadatas, err := pubsignals.ParseQueriesMetadata(
		ctx,
		query.Type,
		ldContextJSON,
		query.CredentialSubject,
		merklizerOptions,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse query metadata: %w", err)
	}

	if len(metadatas) == 0 {
		return nil, fmt.Errorf("no query metadata found")
	}

	// Build circuit query
	metadata := metadatas[0]
	circuitQuery := circuits.Query{
		Operator:  metadata.Operator,
		Values:    metadata.Values,
		SlotIndex: metadata.SlotIndex,
	}

	// Handle merklized credentials
	if metadata.MerklizedSchema {

		merklizer, err := vc.Merklize(ctx, merklize.WithDocumentLoader(documentLoader))
		if err != nil {

			return nil, fmt.Errorf("failed to merklize credential: %w", err)
		}

		if metadata.Path == nil {

			return nil, fmt.Errorf("metadata path is nil for merklized schema")
		}

		proof, value, err := merklizer.Proof(ctx, *metadata.Path)
		if err != nil {

			return nil, fmt.Errorf("failed to get merkle proof: %w", err)
		}

		mtEntry, err := value.MtEntry()
		if err != nil {
			return nil, fmt.Errorf("failed to get mt entry: %w", err)
		}

		circuitQuery.ValueProof = &circuits.ValueProof{
			Path:  metadata.ClaimPathKey,
			Value: mtEntry,
			MTP:   proof,
		}

	}

	// Build circuit inputs
	s.logger.Debug("Building circuit inputs")
	inputs := circuits.AtomicQueryV3Inputs{
		RequestID:                requestId,
		ID:                       &userId,
		ProfileNonce:             big.NewInt(0),
		ClaimSubjectProfileNonce: big.NewInt(0),
		Claim: circuits.ClaimWithSigAndMTPProof{
			IssuerID: &issuerId,
			Claim:    claim,
			NonRevProof: circuits.MTProof{
				Proof: claimNonRevMtp,
				TreeState: circuits.TreeState{
					State:          issuerStateValue,
					ClaimsRoot:     issuerState.ClaimsTree.Root(),
					RevocationRoot: issuerState.RevTree.Root(),
					RootOfRoots:    issuerState.RootsTree.Root(),
				},
			},
			IncProof: &circuits.MTProof{
				Proof: claimIncMtp,
				TreeState: circuits.TreeState{
					State:          issuerStateValue,
					ClaimsRoot:     issuerState.ClaimsTree.Root(),
					RevocationRoot: issuerState.RevTree.Root(),
					RootOfRoots:    issuerState.RootsTree.Root(),
				},
			},
			SignatureProof: &circuits.BJJSignatureProof{
				Signature:       signature,
				IssuerAuthClaim: authClaim,
				IssuerAuthNonRevProof: circuits.MTProof{
					Proof: authClaimNonRevMtp,
					TreeState: circuits.TreeState{
						State:          issuerStateValue,
						ClaimsRoot:     issuerState.ClaimsTree.Root(),
						RevocationRoot: issuerState.RevTree.Root(),
						RootOfRoots:    issuerState.RootsTree.Root(),
					},
				},
				IssuerAuthIncProof: circuits.MTProof{
					Proof: authClaimIncMtp,
					TreeState: circuits.TreeState{
						State:          issuerStateValue,
						ClaimsRoot:     issuerState.ClaimsTree.Root(),
						RevocationRoot: issuerState.RevTree.Root(),
						RootOfRoots:    issuerState.RootsTree.Root(),
					},
				},
			},
		},
		SkipClaimRevocationCheck: query.SkipClaimRevocationCheck,
		Query:                    circuitQuery,
		CurrentTimeStamp:         time.Now().Unix(),
		ProofType:                circuits.ProofType(query.ProofType),
		VerifierID:               &verifierId,
		LinkNonce:                big.NewInt(0),
		NullifierSessionID:       big.NewInt(0),
	}

	return inputs.InputsMarshal()
}

func (s *CircuitService) GenerateSignatureBasedInputs(
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
		RevocationRoot: issuerState.RevTree.Root(),
		RootOfRoots:    issuerState.RootsTree.Root(),
	}

	// 4. get auth claim
	authClaim, err := issuerState.GetAuthClaim()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get auth claim proof
	issuerAuthClaimProof, err := issuerState.GetIncMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. get non-rev auth claim proof
	issuerAuthClaimNonRevProof, err := issuerState.GetNonRevMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 8. get non-rev claim proof
	issuerClaimNonRevProof, err := issuerState.GetNonRevMTProof(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 9. get issuerID
	issuerID := issuerState.GetID()

	inputs := circuits.AtomicQueryV3Inputs{
		RequestID:                requestID,
		ID:                       holderID,
		ProfileNonce:             big.NewInt(0),
		ClaimSubjectProfileNonce: big.NewInt(0),

		Claim: circuits.ClaimWithSigAndMTPProof{
			IssuerID: issuerID,
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

func (s *CircuitService) GenerateMTPBasedInputs(
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
		RevocationRoot: issuerState.RevTree.Root(),
		RootOfRoots:    issuerState.RootsTree.Root(),
	}

	// 3. get claim proof
	issuerClaimProof, err := issuerState.GetIncMTProof(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get non-rev claim proof
	issuerClaimNonRevProof, err := issuerState.GetNonRevMTProof(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get issuer ID
	issuerID := issuerState.GetID()

	inputs := circuits.AtomicQueryV3Inputs{
		RequestID:                requestID,
		ID:                       holderID,
		ProfileNonce:             big.NewInt(0),
		ClaimSubjectProfileNonce: big.NewInt(0),

		Claim: circuits.ClaimWithSigAndMTPProof{
			IssuerID: issuerID,
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

func (s *CircuitService) GenerateAuthV3Inputs(
	ctx context.Context,
	gistRoot *merkletree.Hash,
	gistTree *merkletree.MerkleTree,
	userState *IdentityState,
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
	authClaimProof, err := userState.GetIncMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 4. get non-rev auth claim proof
	authClaimNonRevProof, err := userState.GetNonRevMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get state gist Proof
	gistProof, _, err := gistTree.GenerateProof(ctx, userStateValue.BigInt(), gistTree.Root())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 7. user id
	userID := userState.GetID()

	inputs := circuits.AuthV3Inputs{
		GenesisID:          userID,
		ProfileNonce:       big.NewInt(0),
		AuthClaim:          authClaim,
		AuthClaimIncMtp:    authClaimProof,
		AuthClaimNonRevMtp: authClaimNonRevProof,
		TreeState: circuits.TreeState{
			State:          userStateValue,
			ClaimsRoot:     userState.ClaimsTree.Root(),
			RevocationRoot: userState.RevTree.Root(),
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

func (s *CircuitService) GenerateStateTransitionInputs(
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
	authClaimOldStateProof, err := oldState.GetIncMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 5. get non-rev auth claim proof with old state
	authClaimNonRevOldStateProof, err := oldState.GetNonRevMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 6. get auth claim proof with new state
	authClaimNewStateProof, err := oldState.GetIncMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 8. get user id
	userID := oldState.GetID()

	inputs := circuits.StateTransitionInputs{
		ID: userID,
		OldTreeState: circuits.TreeState{
			State:          oldStateValue,
			ClaimsRoot:     oldState.ClaimsTree.Root(),
			RevocationRoot: oldState.RevTree.Root(),
			RootOfRoots:    oldState.RootsTree.Root(),
		},
		NewTreeState: circuits.TreeState{
			State:          newStateValue,
			ClaimsRoot:     newState.ClaimsTree.Root(),
			RevocationRoot: newState.RevTree.Root(),
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
