package service

import (
	"be/config"
	"be/internal/domain/statistic.go"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
)

type IStatisticService interface {
	GetIssuerStatisticByDID(ctx context.Context, did string) (*dto.IssuerStatisticResponseDto, error)
	GetHolderStatisticByDID(ctx context.Context, did string) (*dto.HolderStatisticResponseDto, error)
	GetVerifierStatisticByDID(ctx context.Context, did string) (*dto.VerifierStatisticResponseDto, error)
}

type StatisticService struct {
	config        *config.Config
	statisticRepo statistic.IStatisticRepository
}

func NewStatisticService(
	config *config.Config,
	statisticRepo statistic.IStatisticRepository,
) IStatisticService {
	return &StatisticService{
		config:        config,
		statisticRepo: statisticRepo,
	}
}

func (s *StatisticService) GetIssuerStatisticByDID(ctx context.Context, did string) (*dto.IssuerStatisticResponseDto, error) {
	entity, err := s.statisticRepo.FindIssuerStatisticByDID(ctx, did)
	if err != nil {
		return nil, &constant.StatisticNotFound
	}
	return &dto.IssuerStatisticResponseDto{
		DocumentNum:          entity.DocumentNum,
		SchemaNum:            entity.SchemaNum,
		CredentialRequestNum: entity.CredentialRequestNum,
		CredentialIssuedNum:  entity.CredentialIssuedNum,
	}, nil
}

func (s *StatisticService) GetHolderStatisticByDID(ctx context.Context, did string) (*dto.HolderStatisticResponseDto, error) {
	entity, err := s.statisticRepo.FindHolderStatisticByDID(ctx, did)
	if err != nil {
		return nil, &constant.StatisticNotFound
	}

	return &dto.HolderStatisticResponseDto{
		CredentialRequestNum:    entity.CredentialRequestNum,
		VerifiableCredentialNum: entity.VerifiableCredentialNum,
		ProofSubmissionNum:      entity.ProofSubmissionNum,
		ProofAcceptedNum:        entity.ProofAcceptedNum,
	}, nil
}

func (s *StatisticService) GetVerifierStatisticByDID(ctx context.Context, did string) (*dto.VerifierStatisticResponseDto, error) {
	entity, err := s.statisticRepo.FindVerifierStatisticByDID(ctx, did)
	if err != nil {
		return nil, &constant.StatisticNotFound
	}

	return &dto.VerifierStatisticResponseDto{
		ProofRequestNum:  entity.ProofRequestNum,
		ProofSubmission:  entity.ProofSubmissionNum,
		ProofAcceptedNum: entity.ProofAcceptedNum,
	}, nil
}
