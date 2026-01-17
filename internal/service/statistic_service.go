package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/proof"
	"be/internal/transport/http/dto"
)

type IStatisticService interface {
}

type StatisticService struct {
	config               *config.Config
	credentialRequest    credential.ICredentialRequestRepository
	proofRequest         proof.IProofRepository
	verifiableCredential credential.IVerifiableCredentialRepository
}

func NewStatisticService(
	config *config.Config,
	credentialRequest credential.ICredentialRequestRepository,
	proofRequest proof.IProofRepository,
	verifiableCredential credential.IVerifiableCredentialRepository) IStatisticService {
	return &StatisticService{
		config:               config,
		credentialRequest:    credentialRequest,
		proofRequest:         proofRequest,
		verifiableCredential: verifiableCredential,
	}
}

func (s *StatisticService) GetIssuerStatistic() *dto.IssuerStatisticResponse {
	return nil
}

func (s *StatisticService) GetHolderStatistic() *dto.HolderStatisticResponse {
	return nil
}

func (s *StatisticService) GetVerifierStatistic() *dto.VerifiableCredentialResponseDto {
	return nil
}
