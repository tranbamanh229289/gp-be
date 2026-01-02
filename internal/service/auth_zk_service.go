package service

import (
	"be/config"
	"be/internal/infrastructure/cache/redis"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iden3/go-circuits/v2"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/iden3comm/v2/protocol"
	"go.uber.org/zap/zapcore"
)

type IAuthZkService interface {
	Register(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error)
	Login(ctx context.Context, authResponse *protocol.AuthorizationResponseMessage) (*dto.IdentityResponseDto, error)
	Challenge(ctx context.Context) (*protocol.AuthorizationRequestMessage, error)
	GetIdentityByRole(ctx context.Context, role string) ([]*dto.IdentityResponseDto, error)
}

type AuthZkService struct {
	config          *config.Config
	logger          *logger.ZapLogger
	redis           *redis.RedisCache
	verifier        *Verifier
	identityService IIdentityService
}

func NewAuthZkService(
	config *config.Config,
	logger *logger.ZapLogger,
	redis *redis.RedisCache,
	identityService IIdentityService,
	verifierService IVerifierService,
) (IAuthZkService, error) {
	return &AuthZkService{
		config:          config,
		logger:          logger,
		redis:           redis,
		verifier:        verifierService.GetVerifier(),
		identityService: identityService,
	}, nil
}

func (s *AuthZkService) Register(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error) {
	return s.identityService.CreateIdentity(ctx, request)

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
	return s.identityService.GetIdentityByDID(ctx, fromDID)
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
	verifierIdentityState, err := s.identityService.GetIdentityState(ctx, verifierPrivateKey.Public())
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

func (s *AuthZkService) GetIdentityByRole(ctx context.Context, role string) ([]*dto.IdentityResponseDto, error) {
	return s.identityService.GetIdentityByRole(ctx, role)
}
