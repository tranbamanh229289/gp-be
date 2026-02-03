package service

import (
	"be/config"
	"be/internal/infrastructure/cache/redis"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iden3/go-circuits/v2"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/iden3comm/v2/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type IAuthZkService interface {
	Register(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error)
	Login(ctx context.Context, authResponse *protocol.AuthorizationResponseMessage) (*dto.ZKLoginResponseDto, error)
	Logout(ctx context.Context, id string) error
	Challenge(ctx context.Context) (*protocol.AuthorizationRequestMessage, error)
	GetIdentityByRole(ctx context.Context, role string) ([]*dto.IdentityResponseDto, error)
	GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error)
	RefreshZKToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponseDto, error)
	VerifyZKToken(tokenString string, tokenType constant.TokenType) (*dto.ZKClaims, error)
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

func (s *AuthZkService) Login(ctx context.Context, authResponse *protocol.AuthorizationResponseMessage) (*dto.ZKLoginResponseDto, error) {
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
	start := time.Now()
	err = s.verifier.VerifyAuthResponse(ctx, *authResponse, authRequest)
	elapsed := time.Since(start)
	fmt.Println("Elapsed ZK Login:", elapsed)
	if err != nil {
		s.logger.Info(fmt.Sprintf("Unauthorize: %s", err))
		return nil, fmt.Errorf("request is unauthorize")
	}
	s.logger.Info("Authorized !")
	fromDID := authResponse.From

	identity, err := s.identityService.GetIdentityByDID(ctx, fromDID)
	if err != nil {
		return nil, fmt.Errorf("get identity error")
	}

	claims := &dto.ZKClaims{
		ID:    identity.PublicID,
		Name:  identity.Name,
		DID:   identity.DID,
		Role:  identity.Role,
		State: identity.State,
	}
	accessToken, err := s.GetZKToken(claims, constant.AccessToken)
	if err != nil {
		s.logger.Info("get access token invalid: ", zap.String("access token:", accessToken), zap.Error(err))
	}
	refreshToken, err := s.GetZKToken(claims, constant.RefreshToken)
	if err != nil {
		s.logger.Error("get refresh token invalid:", zap.String("refresh token:", refreshToken), zap.Error(err))
	}
	publicKey := dto.PublicKeyDto{
		X: identity.PublicKeyX,
		Y: identity.PublicKeyY,
	}

	return &dto.ZKLoginResponseDto{
		Claims:       *claims,
		PublicKey:    publicKey,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
	expiration := 5 * time.Minute
	err = s.redis.Set(ctx, redisKey, redisValue, expiration)

	return &authRequest, nil
}

func (s *AuthZkService) Logout(ctx context.Context, id string) error {
	accessTokenRedisKey := "authzk:" + string(constant.AccessToken) + ":" + id
	refreshTokenRedisKey := "authzk:" + string(constant.RefreshToken) + ":" + id

	err := s.redis.Delete(ctx, accessTokenRedisKey)
	if err != nil {
		return err
	}

	err = s.redis.Delete(ctx, refreshTokenRedisKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthZkService) GetIdentityByRole(ctx context.Context, role string) ([]*dto.IdentityResponseDto, error) {
	return s.identityService.GetIdentityByRole(ctx, role)
}

func (s *AuthZkService) GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error) {
	return s.identityService.GetIdentityByDID(ctx, did)
}

func (s *AuthZkService) RefreshZKToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponseDto, error) {
	claims, err := s.VerifyZKToken(refreshToken, constant.RefreshToken)
	if err != nil {
		return nil, &constant.InvalidToken
	}

	newAccessToken, err := s.GetZKToken(claims, constant.AccessToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate access token %s", err))
		return nil, &constant.InternalServer
	}

	newRefreshToken, err := s.GetZKToken(claims, constant.RefreshToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate refresh token %s", err))
		return nil, &constant.InternalServer
	}

	return &dto.RefreshTokenResponseDto{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthZkService) GetZKToken(claims *dto.ZKClaims, tokenType constant.TokenType) (string, error) {
	now := time.Now().UTC()
	secretKey := []byte(s.config.JWT.Secret)

	var durationTTL time.Duration
	if tokenType == constant.AccessToken {
		durationTTL = s.config.JWT.AccessTokenTTL
	} else {
		durationTTL = s.config.JWT.RefreshTokenTTL
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(durationTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	tokenRedisKey := "authzk:" + string(tokenType) + ":" + claims.ID

	if err := s.redis.Set(context.Background(), tokenRedisKey, tokenString, durationTTL); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthZkService) VerifyZKToken(tokenString string, tokenType constant.TokenType) (*dto.ZKClaims, error) {
	var claims dto.ZKClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	tokenRedisKey := "authzk:" + string(tokenType) + ":" + claims.ID
	tokenRedisValue, err := s.redis.Get(context.Background(), tokenRedisKey).Result()
	if err != nil {
		return nil, err
	}

	if tokenRedisValue != tokenString {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
