package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/proof"
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
)

type IProofService interface {
}

type ProofService struct {
	config       *config.Config
	identityRepo credential.IIdentityRepository
	schemaRepo   schema.ISchemaRepository
	proofRepo    proof.IProofRepository
}

func NewProofService(
	config *config.Config,
	identityRepo credential.IIdentityRepository,
	schemaRepo schema.ISchemaRepository,
	proofRepo proof.IProofRepository,
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

	schema, err := s.schemaRepo.FindSchemaByPublicId(ctx, request.SchemaID)
	if err != nil {
		return nil, err
	}

	proofCreated, err := s.proofRepo.CreateProofRequest(ctx, &proof.ProofRequest{
		VerifierID:        verifier.ID,
		SchemaID:          schema.ID,
		AllowedIssuersDID: request.AllowedIssuersDID,
		Challenge:         request.Challenge,
		Status:            constant.ProofPendingStatus,
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
		AllowedIssuersDID: proofCreated.AllowedIssuersDID,
		Challenge:         proofCreated.Challenge,
		Status:            proofCreated.Status,
		QRCodeData:        proofCreated.QRCodeData,
		CallbackURL:       proofCreated.CallbackURL,
		ExpireAt:          proofCreated.ExpireAt,
	}, nil
}
