package service

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"fmt"

	auth "github.com/iden3/go-iden3-auth/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	"github.com/iden3/iden3comm/v2/protocol"
)

type IAuthZkService interface {
	CreateAuthRequest(ctx context.Context) (*protocol.AuthorizationRequestMessage, error)
	VerifyToken(ctx context.Context, token string, request protocol.AuthorizationRequestMessage) (*protocol.AuthorizationResponseMessage, error)
}

type AuthZkService struct {
	config   *config.Config
	logger   *logger.ZapLogger
	verifier *auth.Verifier
}

func NewAuthZkService(config *config.Config, logger *logger.ZapLogger) IAuthZkService {
	resolvers := map[string]pubsignals.StateResolver{
		config.Blockchain.Resolver: state.NewETHResolver(config.Blockchain.RPC, config.Blockchain.StateContract),
	}
	verifier, err := auth.NewVerifier(loaders.NewEmbeddedKeyLoader(), resolvers)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create verifier %s", err))
		return nil
	}

	return &AuthZkService{
		config:   config,
		logger:   logger,
		verifier: verifier,
	}
}

func (as *AuthZkService) CreateAuthRequest(ctx context.Context) (*protocol.AuthorizationRequestMessage, error) {
	reason := ""
	verifierID := ""
	callbackURL := ""
	request := auth.CreateAuthorizationRequest(reason, verifierID, callbackURL)
	return &request, nil
}

func (as *AuthZkService) CreateAuthRequestWithMessage(ctx context.Context) (*protocol.AuthorizationRequestMessage, error) {
	reason := ""
	message := ""
	verifierID := ""
	callbackURL := ""
	request := auth.CreateAuthorizationRequestWithMessage(reason, message, verifierID, callbackURL)
	return &request, nil
}

func (as *AuthZkService) CreateAuthRequestWithProof(ctx context.Context) (*protocol.AuthorizationRequestMessage, error) {
	reason := ""
	verifierID := ""
	callbackURL := ""
	request := auth.CreateAuthorizationRequest(reason, verifierID, callbackURL)
	var zkProofRequest protocol.ZeroKnowledgeProofRequest
	zkProofRequest.ID = 1
	zkProofRequest.CircuitID = "credentialAtomicQueryV3"
	zkProofRequest.Query = map[string]interface{}{
		"allowedIssuer": []string{"*"},
		"credentialSubject": map[string]interface{}{
			"birthday": map[string]interface{}{
				"$lt": 20020101,
			},
		},
		"context": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v4.jsonld",
		"type":    "KYCAgeCredential",
	}
	request.Body.Scope = append(request.Body.Scope, zkProofRequest)

	return &request, nil
}

func (as *AuthZkService) VerifyToken(ctx context.Context, token string, request protocol.AuthorizationRequestMessage) (*protocol.AuthorizationResponseMessage, error) {
	authResponse, err := as.verifier.FullVerify(ctx, token, request)

	if err != nil {
		return nil, err
	}
	return authResponse, nil
}
