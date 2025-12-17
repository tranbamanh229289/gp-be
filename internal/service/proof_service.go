package service

import (
	"be/config"
	"be/internal/domain/proof"
	"be/internal/infrastructure/database/repository"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
)

type IProofService interface {
}

type ProofService struct {
	config       *config.Config
	identityRepo *repository.IdentityRepository
	schemaRepo   *repository.CredentialRepository
	proofRepo    *repository.ProofRepository
}

func NewProofService(
	config *config.Config,
	identityRepo *repository.IdentityRepository,
	schemaRepo *repository.CredentialRepository,
	proofRepo *repository.ProofRepository,
) IProofService {
	return &ProofService{
		config:       config,
		identityRepo: identityRepo,
		schemaRepo:   schemaRepo,
		proofRepo:    proofRepo,
	}
}

func (s *ProofService) CreateProofRequest(
	ctx context.Context,
	request *dto.ProofRequestCreatedRequestDto,
) (*dto.ProofRequestCreatedResponseDto, error) {
	verifier, err := s.identityRepo.FindIdentityByPublicId(ctx, request.VerifierID)
	if err != nil {
		return nil, err
	}

	schema, err := s.schemaRepo.FindCredentialByPublicId(ctx, request.SchemaID)
	if err != nil {
		return nil, err
	}

	proofCreated, err := s.proofRepo.CreateProofRequest(ctx, &proof.ProofRequest{
		VerifierID:        verifier.ID,
		SchemaID:          schema.ID,
		CircuitType:       request.CircuitType,
		AllowedIssuersDID: request.AllowedIssuersDID,
		Challenge:         request.Challenge,
		Status:            "Pending",
	})
	proofCreated.CallbackURL = ""
	proofCreated.QRCodeData = ""
	if err != nil {
		return nil, &constant.InternalServer
	}

	return &dto.ProofRequestCreatedResponseDto{
		PublicID:          proofCreated.PublicID.String(),
		VerifierID:        verifier.PublicID.String(),
		SchemaID:          schema.PublicID.String(),
		CircuitType:       proofCreated.CircuitType,
		AllowedIssuersDID: proofCreated.AllowedIssuersDID,
		Challenge:         proofCreated.Challenge,
		Status:            proofCreated.Status,
		QRCodeData:        proofCreated.QRCodeData,
		CallbackURL:       proofCreated.CallbackURL,
		ExpireAt:          proofCreated.ExpireAt,
	}, nil
}
