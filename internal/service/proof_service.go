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
	"github.com/iden3/iden3comm/v2/protocol"
	"gorm.io/gorm"
)

type IProofService interface {
	CreateProofRequest(ctx context.Context, request *protocol.AuthorizationRequestMessage) (*dto.ProofRequestResponseDto, error)
	GetProofRequests(ctx context.Context, claims *dto.ZKClaims) ([]*dto.ProofRequestResponseDto, error)
	UpdateProofRequest(ctx context.Context, id string, request *dto.ProofRequestUpdatedRequestDto) error
	VerifyZKProof(ctx context.Context, id string) (*dto.ProofSubmissionResponseDto, error)
	CreateProofSubmission(ctx context.Context, proofSubmission *protocol.AuthorizationResponseMessage) (*dto.ProofSubmissionResponseDto, error)
	GetProofSubmissions(ctx context.Context, claims *dto.ZKClaims) ([]*dto.ProofSubmissionResponseDto, error)
}

type ProofService struct {
	config          *config.Config
	logger          *logger.ZapLogger
	verifier        *Verifier
	identityService IIdentityService
	schemaRepo      schema.ISchemaRepository
	proofRepo       proof.IProofRepository
}

func NewProofService(
	config *config.Config,
	logger *logger.ZapLogger,
	verifierService IVerifierService,
	identityService IIdentityService,
	schemaRepo schema.ISchemaRepository,
	proofRepo proof.IProofRepository,
) IProofService {
	return &ProofService{
		config:          config,
		logger:          logger,
		verifier:        verifierService.GetVerifier(),
		identityService: identityService,
		schemaRepo:      schemaRepo,
		proofRepo:       proofRepo,
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

	verifier, err := s.identityService.GetIdentityByDID(ctx, request.From)
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

func (s *ProofService) VerifyZKProof(ctx context.Context, id string) (*dto.ProofSubmissionResponseDto, error) {
	proofSubmissionEntity, err := s.proofRepo.FindProofSubmissionByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofNotFound
		}
		return nil, &constant.InternalServer
	}

	proofRequestEntity, err := s.proofRepo.FindProofRequestByThreadId(ctx, proofSubmissionEntity.ThreadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofNotFound
		}
		return nil, &constant.InternalServer
	}

	if proofRequestEntity.Status != constant.ProofRequestActiveStatus {
		return nil, errors.New("proof request is not active")
	}

	var proofRequest protocol.AuthorizationRequestMessage = dto.ToAuthorizationRequest(proofRequestEntity)
	var proofSubmission protocol.AuthorizationResponseMessage = dto.ToAuthorizationResponse(proofSubmissionEntity)
	fmt.Println(proofRequest)
	err = s.verifier.VerifyAuthResponse(ctx, proofSubmission, proofRequest)
	var changes map[string]interface{}
	var status constant.ProofSubmissionStatus
	if err != nil {
		s.logger.Info(fmt.Sprintf("Verify Failed: %s", err))
		status = constant.ProofSubmissionFailedStatus
	}
	status = constant.ProofSubmissionSuccessStatus

	changes = map[string]interface{}{"status": status, "verified_date": time.Now().UTC()}
	err = s.proofRepo.UpdateProofSubmission(ctx, proofSubmissionEntity, changes)

	return &dto.ProofSubmissionResponseDto{
		PublicID:     proofSubmissionEntity.PublicID.String(),
		RequestID:    proofRequestEntity.PublicID.String(),
		ThreadID:     proofSubmissionEntity.ThreadID,
		HolderDID:    proofSubmissionEntity.HolderDID,
		HolderName:   proofSubmissionEntity.Holder.Name,
		VerifierDID:  proofRequestEntity.VerifierDID,
		VerifierName: proofRequestEntity.Verifier.Name,
		Message:      proofSubmissionEntity.Message,
		ScopeID:      proofSubmissionEntity.ScopeID,
		CircuitID:    proofSubmissionEntity.CircuitID,
		ZKProof:      proofSubmission.Body.Scope[0].ZKProof,
		CreatedTime:  proofSubmissionEntity.CreatedTime,
		ExpiresTime:  proofSubmissionEntity.ExpiresTime,
		Status:       status,
	}, nil

}

func (s *ProofService) CreateProofSubmission(ctx context.Context, proofSubmission *protocol.AuthorizationResponseMessage) (*dto.ProofSubmissionResponseDto, error) {
	holder, err := s.identityService.GetIdentityByDID(ctx, proofSubmission.From)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	verifier, err := s.identityService.GetIdentityByDID(ctx, proofSubmission.To)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	proofRequestEntity, err := s.proofRepo.FindProofRequestByThreadId(ctx, proofSubmission.ThreadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofNotFound
		}
		return nil, &constant.InternalServer
	}

	scope := proofSubmission.Body.Scope[0]

	zkProofByte, err := json.Marshal(scope.ZKProof)
	if err != nil {
		return nil, err
	}

	proofResponseCreated, err := s.proofRepo.CreateProofSubmission(ctx, &proof.ProofSubmission{
		PublicID:    uuid.New(),
		RequestID:   proofRequestEntity.ID,
		HolderDID:   proofSubmission.From,
		ThreadID:    proofSubmission.ThreadID,
		Message:     proofSubmission.Body.Message,
		ScopeID:     scope.ID,
		CircuitID:   scope.CircuitID,
		ZKProof:     zkProofByte,
		CreatedTime: proofSubmission.CreatedTime,
		ExpiresTime: proofSubmission.ExpiresTime,
		Status:      constant.ProofSubmissionPendingStatus,
	})

	return &dto.ProofSubmissionResponseDto{
		PublicID:     proofResponseCreated.PublicID.String(),
		RequestID:    proofRequestEntity.PublicID.String(),
		ThreadID:     proofResponseCreated.ThreadID,
		HolderDID:    holder.DID,
		HolderName:   holder.Name,
		VerifierDID:  verifier.DID,
		VerifierName: verifier.Name,
		Message:      proofResponseCreated.Message,
		ScopeID:      proofResponseCreated.ScopeID,
		CircuitID:    proofResponseCreated.CircuitID,
		ZKProof:      scope.ZKProof,
		CreatedTime:  proofResponseCreated.CreatedTime,
		ExpiresTime:  proofResponseCreated.ExpiresTime,
		Status:       proofResponseCreated.Status,
	}, nil
}

func (s *ProofService) GetProofSubmissions(ctx context.Context, claims *dto.ZKClaims) ([]*dto.ProofSubmissionResponseDto, error) {
	var (
		proofSubmissions []*proof.ProofSubmission
		err              error
	)
	if claims.Role == constant.IdentityVerifierRole {
		proofSubmissions, err = s.proofRepo.FindAllProofSubmissionsByVerifierDID(ctx, claims.DID)
	} else if claims.Role == constant.IdentityHolderRole {
		proofSubmissions, err = s.proofRepo.FindAllProofSubmissionsByHolderDID(ctx, claims.DID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.ProofNotFound
		}
		return nil, &constant.InternalServer
	}

	var resp []*dto.ProofSubmissionResponseDto
	for _, item := range proofSubmissions {
		resp = append(resp, dto.ToProofSubmissionResponseDto(item))
	}
	return resp, nil
}
