package repository

import (
	"be/internal/domain/proof"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type ProofRepository struct {
	db *postgres.PostgresDB
}

func NewProofRepository(db *postgres.PostgresDB) proof.IProofRepository {
	return &ProofRepository{
		db: db,
	}
}

func (r *ProofRepository) GetProofRequestByPublicId(ctx context.Context, publicId string) (*proof.ProofRequest, error) {
	var entity proof.ProofRequest
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProofRepository) GetProofRequestByRequestId(ctx context.Context, requestId string) (*proof.ProofRequest, error) {
	var entity proof.ProofRequest
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("request_id = ?", requestId).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProofRepository) CreateProofRequest(ctx context.Context, entity *proof.ProofRequest) (*proof.ProofRequest, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ProofRepository) GetProofResponseByPublicId(ctx context.Context, publicId string) (*proof.ProofResponse, error) {
	var entity proof.ProofResponse
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProofRepository) CreateProofResponse(ctx context.Context, entity *proof.ProofResponse) (*proof.ProofResponse, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
