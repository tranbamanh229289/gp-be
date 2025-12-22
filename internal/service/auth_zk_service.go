package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/repository"
	"be/internal/shared/constant"
	"be/internal/shared/types"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/iden3/go-circuits/v2"
	auth "github.com/iden3/go-iden3-auth/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/iden3comm/v2/protocol"
	"gorm.io/gorm"
)

type IAuthZkService interface {
	GetIdentityByPublicId(ctx context.Context, publicId string) (*dto.IdentityResponseDto, error)
	GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error)
	Register(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error)
	Login(ctx context.Context) ([]byte, error)
}

type AuthZkService struct {
	config       *config.Config
	logger       *logger.ZapLogger
	verifier     *auth.Verifier
	identityRepo credential.IIdentityRepository
	mtRepo       repository.IMTRepository
}

func NewAuthZkService(config *config.Config, logger *logger.ZapLogger, identityRepo credential.IIdentityRepository, mtRepo repository.IMTRepository) (IAuthZkService, error) {
	resolvers := map[string]pubsignals.StateResolver{
		config.Blockchain.Resolver: state.NewETHResolver(config.Blockchain.RPC, config.Blockchain.StateContract),
	}
	keyLoader := loaders.FSKeyLoader{Dir: config.Circuit.AuthV3VerifyingKeyPath}
	verifier, err := auth.NewVerifier(keyLoader, resolvers)
	if err != nil {
		return nil, fmt.Errorf("failed to create verifier %s", err)
	}

	return &AuthZkService{
		config:       config,
		logger:       logger,
		verifier:     verifier,
		identityRepo: identityRepo,
		mtRepo:       mtRepo,
	}, nil
}

func (s *AuthZkService) GetIdentityByPublicId(ctx context.Context, publicId string) (*dto.IdentityResponseDto, error) {
	identity, err := s.identityRepo.FindIdentityByPublicId(ctx, publicId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identity), nil
}

func (s *AuthZkService) GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error) {
	identity, err := s.identityRepo.FindIdentityByDID(ctx, did)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identity), nil
}

func (s *AuthZkService) Register(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error) {
	X := new(big.Int)
	Y := new(big.Int)
	X.SetString(request.PublicKeyX, 10)
	Y.SetString(request.PublicKeyY, 10)
	publicKey := &babyjub.PublicKey{X: X, Y: Y}

	genesisState, err := s.CreateGenesisState(ctx, publicKey)
	if err != nil {
		return nil, err
	}
	stateValue, err := genesisState.GetStateValue()
	if err != nil {
		return nil, err
	}
	did, err := genesisState.GetDID()
	if err != nil {
		return nil, err
	}
	identityEntity := &credential.Identity{
		PublicID:   uuid.New(),
		DID:        types.DID(did.String()),
		State:      types.Hash(stateValue.String()),
		Name:       request.Name,
		Role:       request.Role,
		ClaimsMTID: genesisState.ClaimsMTID,
		RevMTID:    genesisState.RevMTID,
		RootsMTID:  genesisState.RootsMTID,
		PublicKeyX: publicKey.X.String(),
		PublicKeyY: publicKey.Y.String(),
	}

	identityCreated, err := s.identityRepo.CreateIdentity(ctx, identityEntity)
	if err != nil {
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identityCreated), nil
}

func (s *AuthZkService) Login(ctx context.Context) ([]byte, error) {
	verifierPrivateKeyBytes, err := hex.DecodeString(s.config.Iden3.VerifierPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed decode private key %s", err)
	}
	if len(verifierPrivateKeyBytes) != 32 {
		return nil, fmt.Errorf("invalid private key length: %d", len(verifierPrivateKeyBytes))
	}
	verifierPrivateKey := babyjub.PrivateKey(verifierPrivateKeyBytes)
	verifierIdentity, err := s.GetIdentityState(ctx, verifierPrivateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("failed to load verifier wallet %s", err)
	}
	verifierDID, _ := verifierIdentity.GetDID()

	baseURL := s.config.GetBaseURL()
	callbackURL := baseURL + "/api/v1/auth/callback"
	challenge, _ := s.generateRandomChallenge()
	reason := "Authenticate"
	message := "Please sign in with zk proof"

	var scopes []protocol.ZeroKnowledgeProofRequest
	scopes = append(scopes, protocol.ZeroKnowledgeProofRequest{
		ID:        uuid.New().ID(),
		CircuitID: string(circuits.AuthV3CircuitID),
		Params: map[string]interface{}{
			"challenge": challenge,
		},
	})
	request := auth.CreateAuthorizationRequestWithMessage(reason, message, verifierDID.String(), callbackURL)
	request.Body.Scope = scopes

	qrCode, err := json.Marshal(request)

	return qrCode, nil
}

func (s *AuthZkService) Callback(ctx context.Context) {

}

func (s *AuthZkService) generateRandomChallenge() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthZkService) CreateGenesisState(ctx context.Context, publicKey *babyjub.PublicKey) (*IdentityState, error) {
	// create auth claim
	authSchemaHash, _ := core.NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	revNonce := uint64(1)
	authClaim, _ := core.NewClaim(authSchemaHash, core.WithIndexDataInts(publicKey.X, publicKey.Y), core.WithRevocationNonce(revNonce))

	// add auth claim into claimsTree
	claimsTree, claimsMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	authHi, authHv, _ := authClaim.HiHv()
	claimsTree.Add(ctx, authHi, authHv)

	// create empty revocationsTree
	revTree, revMTID, _ := s.mtRepo.NewMerkleTree(ctx)

	// create empty rootsTree
	rootsTree, rootsMTID, _ := s.mtRepo.NewMerkleTree(ctx)

	identity := &IdentityState{
		ClaimsTree:      claimsTree,
		ClaimsMTID:      claimsMTID,
		RevocationsTree: revTree,
		RevMTID:         revMTID,
		RootsTree:       rootsTree,
		RootsMTID:       rootsMTID,
		PublicKey:       publicKey,
	}

	return identity, nil
}

func (s *AuthZkService) GetIdentityState(ctx context.Context, publicKey *babyjub.PublicKey) (*IdentityState, error) {
	identity, err := s.identityRepo.FindIdentityByPublicKey(ctx, publicKey.X.String(), publicKey.Y.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	claimsTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.ClaimsMTID)
	revTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.RevMTID)
	rootsTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.RootsMTID)

	identityState := &IdentityState{
		ClaimsTree:      claimsTree,
		ClaimsMTID:      identity.ClaimsMTID,
		RevocationsTree: revTree,
		RevMTID:         identity.RevMTID,
		RootsTree:       rootsTree,
		RootsMTID:       identity.RootsMTID,
		PublicKey:       publicKey,
	}

	return identityState, nil
}
