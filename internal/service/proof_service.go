package service

import (
	"be/config"
	"be/internal/domain/proof"
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/iden3comm/v2/packers"
	"github.com/iden3/iden3comm/v2/protocol"
	"gorm.io/gorm"
)

type IProofService interface {
	CreateProofRequest(ctx context.Context, request *protocol.AuthorizationRequestMessage) (*dto.ProofRequestResponseDto, error)
	GetProofRequests(ctx context.Context, claims *dto.ZKClaims) ([]*dto.ProofRequestResponseDto, error)
	UpdateProofRequest(ctx context.Context, id string, request *dto.ProofRequestUpdatedRequestDto) error
	CreateProofResponse(ctx context.Context, proofResponse *protocol.AuthorizationResponseMessage) (*dto.ProofVerificationResponseDto, error)
	GetProofResponses(ctx context.Context) ([]*dto.ProofVerificationResponseDto, error)
}

type ProofService struct {
	config       *config.Config
	logger       *logger.ZapLogger
	verifier     *Verifier
	identityRepo schema.IIdentityRepository
	schemaRepo   schema.ISchemaRepository
	proofRepo    proof.IProofRepository
}

func NewProofService(
	config *config.Config,
	logger *logger.ZapLogger,
	verifierService IVerifierService,
	identityRepo schema.IIdentityRepository,
	schemaRepo schema.ISchemaRepository,
	proofRepo proof.IProofRepository,
) IProofService {
	return &ProofService{
		config:       config,
		logger:       logger,
		verifier:     verifierService.GetVerifier(),
		identityRepo: identityRepo,
		schemaRepo:   schemaRepo,
		proofRepo:    proofRepo,
	}
}

func (s *ProofService) CreateProofRequest(
	ctx context.Context,
	request *protocol.AuthorizationRequestMessage,
) (*dto.ProofRequestResponseDto, error) {
	scopes := request.Body.Scope
	if len(scopes) != 1 {
		return nil, errors.New("required len scope = 1")
	}

	var proofQuery pubsignals.Query
	queryBytes, err := json.Marshal(scopes[0].Query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}
	if err := json.Unmarshal(queryBytes, &proofQuery); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query: %w", err)
	}

	nullifierSessionId, ok := scopes[0].Params["nullifierSessionId"].(string)
	if !ok {
		nullifierSessionId = ""
	}

	schema, err := s.schemaRepo.FindSchemaByContextURL(ctx, proofQuery.Context)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}

	if proofQuery.Type != schema.Type {
		return nil, errors.New("schema type not found")
	}

	verifier, err := s.identityRepo.FindIdentityByDID(ctx, request.From)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	entity := &proof.ProofRequest{
		PublicID:                 uuid.New(),
		ThreadID:                 request.ThreadID,
		VerifierDID:              request.From,
		CallbackURL:              request.Body.CallbackURL,
		Reason:                   request.Body.Reason,
		Message:                  request.Body.Message,
		ScopeID:                  scopes[0].ID,
		CircuitID:                scopes[0].CircuitID,
		AllowedIssuers:           proofQuery.AllowedIssuers,
		SchemaID:                 schema.ID,
		CredentialSubject:        proofQuery.CredentialSubject,
		ProofType:                proofQuery.ProofType,
		SkipClaimRevocationCheck: proofQuery.SkipClaimRevocationCheck,
		GroupID:                  proofQuery.GroupID,
		NullifierSession:         nullifierSessionId,
		Status:                   constant.ProofRequestActiveStatus,
		ExpiresTime:              request.ExpiresTime,
		CreatedTime:              request.CreatedTime,
	}
	proofCreated, err := s.proofRepo.CreateProofRequest(ctx, entity)

	if err != nil {
		return nil, &constant.InternalServer
	}

	return &dto.ProofRequestResponseDto{
		PublicID:                 proofCreated.PublicID.String(),
		ThreadID:                 proofCreated.ThreadID,
		VerifierDID:              proofCreated.VerifierDID,
		VerifierName:             verifier.Name,
		CallbackURL:              proofCreated.CallbackURL,
		Reason:                   proofCreated.Reason,
		Message:                  proofCreated.Message,
		ScopeID:                  proofCreated.ScopeID,
		CircuitID:                proofCreated.CircuitID,
		AllowedIssuers:           proofCreated.AllowedIssuers,
		CredentialSubject:        proofCreated.CredentialSubject,
		SchemaID:                 schema.PublicID.String(),
		Context:                  schema.ContextURL,
		Type:                     schema.Type,
		ProofType:                proofCreated.ProofType,
		SkipClaimRevocationCheck: proofCreated.SkipClaimRevocationCheck,
		GroupID:                  proofCreated.GroupID,
		NullifierSession:         proofCreated.NullifierSession,
		Status:                   proofCreated.Status,
		ExpiresTime:              proofCreated.ExpiresTime,
		CreatedTime:              proofCreated.CreatedTime,
	}, nil
}

func (s *ProofService) GetProofRequests(ctx context.Context, claims *dto.ZKClaims) ([]*dto.ProofRequestResponseDto, error) {
	var (
		proofRequests []*proof.ProofRequest
		err           error
	)
	if claims.Role != constant.IdentityIssuerRole {
		proofRequests, err = s.proofRepo.FindAllProofRequests(ctx)
	} else {
		proofRequests, err = s.proofRepo.FindAllProofRequestsByVerifierDID(ctx, claims.DID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofNotFound
		}
		return nil, &constant.InternalServer
	}

	var resp []*dto.ProofRequestResponseDto
	for _, item := range proofRequests {
		resp = append(resp, dto.ToProofRequestResponseDto(item))
	}
	return resp, nil
}

func (s *ProofService) UpdateProofRequest(ctx context.Context, id string, request *dto.ProofRequestUpdatedRequestDto) error {
	entity, err := s.proofRepo.FindProofRequestByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.ProofNotFound
		}
		return &constant.InternalServer
	}

	changes := map[string]interface{}{"status": request.Status}
	return s.proofRepo.UpdateProofRequest(ctx, entity, changes)
}

func (s *ProofService) CreateProofResponse(ctx context.Context, proofResponse *protocol.AuthorizationResponseMessage) (*dto.ProofVerificationResponseDto, error) {

	proofRequestEntity, err := s.proofRepo.FindProofRequestByThreadId(ctx, proofResponse.ThreadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("proof request not found")
		}
		return nil, &constant.InternalServer
	}

	if proofRequestEntity.Status != constant.ProofRequestActiveStatus {
		return nil, errors.New("proof request is not active")
	}

	if proofRequestEntity.ExpiresTime != nil && time.Now().Unix() > *proofRequestEntity.ExpiresTime {
		changes := map[string]interface{}{"status": constant.ProofRequestExpiredStatus}
		_ = s.proofRepo.UpdateProofRequest(ctx, proofRequestEntity, changes)
		return nil, errors.New("proof request has expired")
	}

	holderDID := proofResponse.From
	_, err = s.identityRepo.FindIdentityByDID(ctx, holderDID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prover identity not found")
		}
		return nil, &constant.InternalServer
	}

	if proofResponse.To != proofRequestEntity.VerifierDID {
		return nil, errors.New("verifier DID mismatch")
	}
	var proofRequest *protocol.AuthorizationRequestMessage = s.toAuthorizationRequest(proofRequestEntity)
	err = s.verifier.VerifyAuthResponse(ctx, *proofResponse, *proofRequest)
	if err != nil {
		s.logger.Info(fmt.Sprintf("Verify Failed: %s", err))
	}

	proofResponseCreated, err := s.proofRepo.CreateProofResponse(ctx, &proof.ProofResponse{
		PublicID:  uuid.New(),
		RequestID: proofRequestEntity.ID,
		HolderDID: holderDID,
		Status:    constant.ProofResponseSuccessStatus,
	})

	return &dto.ProofVerificationResponseDto{
		Status:     constant.ProofResponseSuccessStatus,
		HolderDID:  holderDID,
		ThreadID:   proofRequest.ThreadID,
		VerifiedAt: proofResponseCreated.CreatedAt,
	}, nil

}

func (s *ProofService) GetProofResponses(ctx context.Context) ([]*dto.ProofVerificationResponseDto, error) {
	proofResponses, err := s.proofRepo.FindAllProofResponses(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofNotFound
		}
		return nil, &constant.InternalServer
	}

	var resp []*dto.ProofVerificationResponseDto
	for _, item := range proofResponses {
		resp = append(resp, &dto.ProofVerificationResponseDto{
			Status:     item.Status,
			HolderDID:  item.HolderDID,
			ThreadID:   item.ProofRequest.ThreadID,
			VerifiedAt: item.CreatedAt,
		})
	}
	return resp, nil
}

func (s *ProofService) toAuthorizationRequest(pr *proof.ProofRequest) *protocol.AuthorizationRequestMessage {
	query := make(map[string]interface{})
	params := make(map[string]interface{})
	query["allowedIssuers"] = pr.AllowedIssuers
	query["context"] = pr.Schema.ContextURL
	query["type"] = pr.Schema.Type
	query["credentialSubject"] = pr.CredentialSubject
	query["proofType"] = pr.ProofType
	query["skipClaimRevocationCheck"] = pr.SkipClaimRevocationCheck
	query["groupId"] = pr.GroupID
	params["nullifierSessionId"] = pr.NullifierSession

	return &protocol.AuthorizationRequestMessage{
		ID:       pr.ThreadID,
		From:     pr.VerifierDID,
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.AuthorizationRequestMessageType,
		ThreadID: pr.ThreadID,
		Body: protocol.AuthorizationRequestMessageBody{
			CallbackURL: pr.CallbackURL,
			Reason:      pr.Reason,
			Message:     pr.Message,
			Scope: []protocol.ZeroKnowledgeProofRequest{
				{
					ID:        pr.ScopeID,
					CircuitID: pr.CircuitID,
					Query:     query,
					Params:    params,
				},
			},
		},

		ExpiresTime: pr.ExpiresTime,
		CreatedTime: pr.CreatedTime,
	}
}
