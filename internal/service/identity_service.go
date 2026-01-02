package service

import (
	"be/config"
	"be/internal/domain/schema"
	"be/internal/infrastructure/database/repository"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
	"errors"
	"math/big"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"gorm.io/gorm"
)

type IIdentityService interface {
	GetIdentityByPublicId(ctx context.Context, publicId string) (*dto.IdentityResponseDto, error)
	GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error)
	GetIdentityByRole(ctx context.Context, role string) ([]*dto.IdentityResponseDto, error)
	CreateIdentity(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error)

	GetIdentityState(ctx context.Context, publicKey *babyjub.PublicKey) (*IdentityState, error)
	GetIdentityStateByDID(ctx context.Context, didStr string) (*IdentityState, error)
}

type IdentityService struct {
	config       *config.Config
	identityRepo schema.IIdentityRepository
	mtRepo       repository.IMTRepository
}

func NewIdentityService(config *config.Config, identityRepo schema.IIdentityRepository, mtRepo repository.IMTRepository) IIdentityService {
	return &IdentityService{
		config:       config,
		identityRepo: identityRepo,
		mtRepo:       mtRepo,
	}
}

func (s *IdentityService) GetIdentityByPublicId(ctx context.Context, publicId string) (*dto.IdentityResponseDto, error) {
	identity, err := s.identityRepo.FindIdentityByPublicId(ctx, publicId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identity), nil
}

func (s *IdentityService) GetIdentityByDID(ctx context.Context, did string) (*dto.IdentityResponseDto, error) {
	identity, err := s.identityRepo.FindIdentityByDID(ctx, did)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identity), nil
}

func (s *IdentityService) GetIdentityByRole(ctx context.Context, role string) ([]*dto.IdentityResponseDto, error) {
	identities, err := s.identityRepo.FindIdentityByRole(ctx, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	var resp []*dto.IdentityResponseDto
	for _, item := range identities {
		resp = append(resp, dto.ToIdentityResponseDto(item))
	}
	return resp, nil
}

func (s *IdentityService) CreateIdentity(ctx context.Context, request *dto.IdentityCreatedRequestDto) (*dto.IdentityResponseDto, error) {
	X := new(big.Int)
	Y := new(big.Int)
	X.SetString(request.PublicKeyX, 10)
	Y.SetString(request.PublicKeyY, 10)
	publicKey := &babyjub.PublicKey{X: X, Y: Y}

	state, err := s.CreateIdentityState(ctx, publicKey)
	if err != nil {
		return nil, err
	}
	stateValue, err := state.GetStateValue()
	if err != nil {
		return nil, err
	}
	did := state.GetDID()
	identityEntity := &schema.Identity{
		PublicID:   uuid.New(),
		DID:        did.String(),
		State:      stateValue.Hex(),
		Name:       request.Name,
		Role:       request.Role,
		ClaimsMTID: state.ClaimsMTID,
		RevMTID:    state.RevMTID,
		RootsMTID:  state.RootsMTID,
		PublicKeyX: publicKey.X.String(),
		PublicKeyY: publicKey.Y.String(),
	}

	identityCreated, err := s.identityRepo.CreateIdentity(ctx, identityEntity)
	if err != nil {
		return nil, &constant.InternalServer
	}
	return dto.ToIdentityResponseDto(identityCreated), nil
}

func (s *IdentityService) CreateIdentityState(ctx context.Context, publicKey *babyjub.PublicKey) (*IdentityState, error) {
	claimsTree, claimsMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	revTree, revMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	rootsTree, rootsMTID, _ := s.mtRepo.NewMerkleTree(ctx)
	identityState := &IdentityState{
		PublicKey:  publicKey,
		ClaimsTree: claimsTree,
		ClaimsMTID: claimsMTID,
		RevTree:    revTree,
		RevMTID:    revMTID,
		RootsTree:  rootsTree,
		RootsMTID:  rootsMTID,
	}
	authClaim, _ := identityState.GetAuthClaim()
	err := identityState.AddClaim(ctx, authClaim)
	if err != nil {
		return nil, err
	}
	state, _ := identityState.GetStateValue()
	typ, _ := core.BuildDIDType(core.DIDMethodPolygonID, core.Polygon, core.Amoy)
	did, _ := core.NewDIDFromIdenState(typ, state.BigInt())
	identityState.DID = did

	return identityState, nil
}

func (s *IdentityService) UpdateState(ctx context.Context, id, state string) error {
	identity, err := s.identityRepo.FindIdentityByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.IdentityNotFound
		}
		return &constant.InternalServer
	}
	identity.State = state
	changes := map[string]interface{}{"state": state}
	return s.identityRepo.UpdateIdentity(ctx, identity, changes)
}

func (s *IdentityService) GetIdentityState(ctx context.Context, publicKey *babyjub.PublicKey) (*IdentityState, error) {
	identity, err := s.identityRepo.FindIdentityByPublicKey(ctx, publicKey.X.String(), publicKey.Y.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	claimsTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.ClaimsMTID)
	revTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.RevMTID)
	rootsTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.RootsMTID)
	did, _ := w3c.ParseDID(identity.DID)
	identityState := &IdentityState{
		PublicKey:  publicKey,
		DID:        did,
		ClaimsTree: claimsTree,
		ClaimsMTID: identity.ClaimsMTID,
		RevTree:    revTree,
		RevMTID:    identity.RevMTID,
		RootsTree:  rootsTree,
		RootsMTID:  identity.RootsMTID,
	}

	return identityState, nil
}

func (s *IdentityService) GetIdentityStateByDID(ctx context.Context, didStr string) (*IdentityState, error) {
	identity, err := s.identityRepo.FindIdentityByDID(ctx, didStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	X := new(big.Int)
	Y := new(big.Int)
	X.SetString(identity.PublicKeyX, 10)
	Y.SetString(identity.PublicKeyY, 10)

	publicKey := &babyjub.PublicKey{X: X, Y: Y}
	claimsTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.ClaimsMTID)
	revTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.RevMTID)
	rootsTree, _ := s.mtRepo.LoadMerkleTree(ctx, identity.RootsMTID)
	did, _ := w3c.ParseDID(identity.DID)
	identityState := &IdentityState{
		PublicKey:  publicKey,
		DID:        did,
		ClaimsTree: claimsTree,
		ClaimsMTID: identity.ClaimsMTID,
		RevTree:    revTree,
		RevMTID:    identity.RevMTID,
		RootsTree:  rootsTree,
		RootsMTID:  identity.RootsMTID,
	}

	return identityState, nil
}
