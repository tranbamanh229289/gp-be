package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICredentialService interface {
}
type CredentialService struct {
	config                *config.Config
	credentialRequestRepo credential.ICredentialRequestRepository
	identityRepo          credential.IIdentityRepository
	schemaRepo            schema.ISchemaRepository
}

func NewCredentialService(
	config *config.Config,
	credentialRequestRepo credential.ICredentialRequestRepository,
	identityRepo credential.IIdentityRepository,
	schemaRepo schema.ISchemaRepository,
) ICredentialService {
	return &CredentialService{
		config:                config,
		credentialRequestRepo: credentialRequestRepo,
		identityRepo:          identityRepo,
		schemaRepo:            schemaRepo,
	}
}

func (r *CredentialService) CreateCredentialRequest(ctx context.Context, request *dto.CredentialRequestCreatedRequestDto) (*dto.CredentialRequestResponseDto, error) {
	holder, err := r.identityRepo.FindIdentityByPublicId(ctx, request.HolderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	schema, err := r.schemaRepo.FindSchemaByPublicId(ctx, request.SchemaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}

	credentialRequest, err := r.credentialRequestRepo.CreateCredentialRequest(ctx, &credential.CredentialRequest{
		PublicID:       uuid.New(),
		HolderID:       holder.ID,
		SchemaID:       schema.ID,
		CredentialData: request.CredentialData,
		Status:         constant.CredentialRequestRejectedStatus,
	})
	if err != nil {
		return nil, &constant.InternalServer
	}

	return &dto.CredentialRequestResponseDto{
		PublicID:       credentialRequest.PublicID.String(),
		HolderID:       holder.PublicID.String(),
		SchemaID:       schema.PublicID.String(),
		CredentialData: credentialRequest.CredentialData,
		Status:         credentialRequest.Status,
	}, nil
}

func (r *CredentialService) GetCredentialRequests(ctx context.Context) ([]*dto.CredentialRequestResponseDto, error) {
	credentialRequests, err := r.credentialRequestRepo.FindAllCredentialRequest(ctx)
	if err != nil {
		return nil, &constant.InternalServer
	}
	var credentialRequestResponseDto []*dto.CredentialRequestResponseDto

	for _, c := range credentialRequests {
		credentialRequestResponseDto = append(credentialRequestResponseDto, &dto.CredentialRequestResponseDto{
			PublicID:       c.PublicID.String(),
			HolderID:       c.Holder.PublicID.String(),
			SchemaID:       c.Schema.PublicID.String(),
			CredentialData: c.CredentialData,
			Status:         c.Status,
		})
	}
	return credentialRequestResponseDto, nil
}

func (r *CredentialService) IssueCredential(ctx context.Context, requestID string) (*credential.Credential, error) {
	credentialRequest, err := r.credentialRequestRepo.FindCredentialRequestByPublicId(ctx, requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.CredentialRequestNotFound
		}
		return nil, &constant.InternalServer
	}
	changes := map[string]interface{}{"status": constant.CredentialRequestApprovedStatus}
	err = r.credentialRequestRepo.UpdateCredentialRequest(ctx, credentialRequest, changes)
	if err != nil {
		return nil, &constant.InternalServer
	}
	return nil, nil
}
