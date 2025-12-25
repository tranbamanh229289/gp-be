package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/repository"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/iden3/go-circuits/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/iden3comm/v2/protocol"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type IAuthZkService interface {
	GetIdentityByPublicId(ctx context.Context, publicId string) (*dto.IdentityResponseDto, error)
	GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error)
	Register(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error)
	Login(ctx context.Context, authResponse *protocol.AuthorizationResponseMessage) (*dto.IdentityResponseDto, error)
	Challenge(ctx context.Context) (*protocol.AuthorizationRequestMessage, error)
}

type AuthZkService struct {
	config       *config.Config
	logger       *logger.ZapLogger
	redis        *redis.RedisCache
	verifier     *Verifier
	identityRepo credential.IIdentityRepository
	mtRepo       repository.IMTRepository
}

func NewAuthZkService(config *config.Config, logger *logger.ZapLogger, redis *redis.RedisCache, identityRepo credential.IIdentityRepository, mtRepo repository.IMTRepository) (IAuthZkService, error) {
	resolvers := map[string]pubsignals.StateResolver{
		config.Blockchain.Resolver: state.NewETHResolver(config.Blockchain.RPC, config.Blockchain.StateContract),
	}

	dir := config.Circuit.VerifyingKey
	keyLoader := loaders.FSKeyLoader{Dir: dir}
	verifier, err := NewVerifier(keyLoader, resolvers)
	if err != nil {
		return nil, fmt.Errorf("failed to create verifier %s", err)
	}

	return &AuthZkService{
		config:       config,
		logger:       logger,
		redis:        redis,
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

	state, err := s.CreateIdentityState(ctx, publicKey)
	if err != nil {
		return nil, err
	}
	stateValue, err := state.GetStateValue()
	if err != nil {
		return nil, err
	}
	did := state.GetDID()
	identityEntity := &credential.Identity{
		PublicID:   uuid.New(),
		DID:        did.String(),
		State:      stateValue.Hex(),
		Name:       request.Name,
		Role:       request.Role,
		ClaimsMTID: state.ClaimsMTID,
		RevMTID:    state.RevMTID,
		RootsMTID:  state.RootsMTID,
		PublicKeyX: publicKey.X.String(),
		PublicKeyY: publicKey.Y.String(),
	}

	identityCreated, err := s.identityRepo.CreateIdentity(ctx, identityEntity)
	if err != nil {
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identityCreated), nil
}

func (s *AuthZkService) Login(ctx context.Context, authResponse *protocol.AuthorizationResponseMessage) (*dto.IdentityResponseDto, error) {
	redisKey := fmt.Sprintf("authzk:request:id:%s", authResponse.ID)
	redisValue, err := s.redis.Get(ctx, redisKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get value redis with key: %s", redisKey)
	}
	var authRequest protocol.AuthorizationRequestMessage
	err = json.Unmarshal(redisValue, &authRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth request")
	}

	err = s.verifier.VerifyAuthResponse(ctx, *authResponse, authRequest)
	if err != nil {
		s.logger.Info(fmt.Sprintf("Unauthorize: %s", err))
		return nil, fmt.Errorf("request is unauthorize")
	}
	s.logger.Info("Authorized !")

	fromDID := authResponse.From
	identity, err := s.identityRepo.FindIdentityByDID(ctx, fromDID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.ToIdentityResponseDto(identity), nil
}

func (s *AuthZkService) Challenge(ctx context.Context) (*protocol.AuthorizationRequestMessage, error) {
	verifierPrivateKeyBytes, err := hex.DecodeString(s.config.Iden3.VerifierPrivateKey)
	if err != nil {
		s.logger.Info("error %w", zapcore.Field{String: err.Error()})
		return nil, fmt.Errorf("failed decode private key %s", err)
	}
	if len(verifierPrivateKeyBytes) != 32 {
		s.logger.Info("private key > 32 char")
		return nil, fmt.Errorf("invalid private key length: %d", len(verifierPrivateKeyBytes))
	}
	verifierPrivateKey := babyjub.PrivateKey(verifierPrivateKeyBytes)
	verifierIdentityState, err := s.GetIdentityState(ctx, verifierPrivateKey.Public())
	if err != nil {
		s.logger.Info("error %w", zapcore.Field{String: err.Error()})
		return nil, fmt.Errorf("failed to load verifier wallet %s", err)
	}
	verifierDID := verifierIdentityState.GetDID()

	callbackURL := "authzk/login"
	challenge, _ := rand.Prime(rand.Reader, 32)
	reason := "Authenticate"
	message := "Please sign in with zk proof"
	var scopes []protocol.ZeroKnowledgeProofRequest
	scopes = append(scopes, protocol.ZeroKnowledgeProofRequest{
		ID:        uint32(challenge.Uint64()),
		CircuitID: string(circuits.AuthV3CircuitID),
		Params: map[string]interface{}{
			"challenge": challenge,
		},
	})

	authRequest := CreateAuthorizationRequestWithMessage(reason, message, verifierDID.String(), callbackURL)
	authRequest.Body.Scope = scopes

	redisKey := fmt.Sprintf("authzk:request:id:%s", authRequest.ID)
	redisValue, err := json.Marshal(authRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth request: %w", err)
	}
	expiration := 2 * time.Minute
	err = s.redis.Set(ctx, redisKey, redisValue, expiration)

	return &authRequest, nil
}

func (s *AuthZkService) CreateIdentityState(ctx context.Context, publicKey *babyjub.PublicKey) (*IdentityState, error) {
	claimsTree, claimsMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	revTree, revMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	rootsTree, rootsMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	identityState := &IdentityState{
		PublicKey:  publicKey,
		ClaimsTree: claimsTree,
		ClaimsMTID: claimsMTID,
		RevTree:    revTree,
		RevMTID:    revMTID,
		RootsTree:  rootsTree,
		RootsMTID:  rootsMTID,
	}
	authClaim, _ := identityState.GetAuthClaim()
	err := identityState.AddClaim(ctx, authClaim)
	if err != nil {
		return nil, err
	}
	state, _ := identityState.GetStateValue()
	typ, _ := core.BuildDIDType(core.DIDMethodPolygonID, core.Polygon, core.Amoy)
	did, _ := core.NewDIDFromIdenState(typ, state.BigInt())
	identityState.DID = did

	return identityState, nil
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
	did, _ := w3c.ParseDID(identity.DID)
	identityState := &IdentityState{
		PublicKey:  publicKey,
		DID:        did,
		ClaimsTree: claimsTree,
		ClaimsMTID: identity.ClaimsMTID,
		RevTree:    revTree,
		RevMTID:    identity.RevMTID,
		RootsTree:  rootsTree,
		RootsMTID:  identity.RootsMTID,
	}

	return identityState, nil
}
