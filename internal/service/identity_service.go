package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/repository"
	"be/internal/infrastructure/message_queue/kafka"
	"be/internal/shared/constant"
	"be/internal/shared/types"
	"be/internal/transport/http/dto"
	"context"
	"math/big"

	"github.com/google/uuid"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type IHolderService interface {
	CreateHolderGenesisState(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityCreatedResponseDto, error)
}

type HolderService struct {
	config       *config.Config
	claimService *ClaimService
	identityRepo *repository.IdentityRepository
	kafkaManager *kafka.Manager
}

func NewHolderService(
	config *config.Config,
	kafkaManager *kafka.Manager,
	identityRepo *repository.IdentityRepository,

) IHolderService {
	kafkaManager.CreateProducer("holder-service")
	return &HolderService{
		config:       config,
		kafkaManager: kafkaManager,
		identityRepo: identityRepo,
	}
}

func (s *HolderService) CreateHolderGenesisState(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityCreatedResponseDto, error) {
	X := new(big.Int)
	Y := new(big.Int)
	X.SetString(request.PublicKeyX, 10)
	Y.SetString(request.PublicKeyY, 10)
	publicKey := &babyjub.PublicKey{X: X, Y: Y}

	genesisState := s.claimService.CreateGenesisState(ctx, publicKey)
	identityEntity := &credential.Identity{
		PublicID:   uuid.New(),
		DID:        types.DID(genesisState.DID.String()),
		State:      types.Hash(genesisState.State.String()),
		Type:       "holder",
		PublicKeyX: publicKey.X.String(),
		PublicKeyY: publicKey.Y.String(),
	}
	identityCreated, err := s.identityRepo.CreateIdentity(ctx, identityEntity)
	if err != nil {
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityCreatedResponseDto(identityCreated), nil
}
