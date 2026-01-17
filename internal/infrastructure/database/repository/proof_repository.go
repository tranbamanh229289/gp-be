package repository

import (
	"be/internal/domain/proof"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
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

func (r *ProofRepository) FindProofRequestByPublicId(ctx context.Context, publicId string) (*proof.ProofRequest, error) {
	var entity proof.ProofRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Verifier").Preload("Schema").Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProofRepository) FindProofRequestByThreadId(ctx context.Context, threadId string) (*proof.ProofRequest, error) {
	var entity proof.ProofRequest
	if err := r.db.GetGormDB().WithContext(ctx).Where("thread_id = ?", threadId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProofRepository) FindAllProofRequests(ctx context.Context) ([]*proof.ProofRequest, error) {
	var entity []*proof.ProofRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Verifier").Preload("Schema").Find(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
func (r *ProofRepository) FindAllProofRequestsByVerifierDID(ctx context.Context, did string) ([]*proof.ProofRequest, error) {
	var entity []*proof.ProofRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Verifier").Preload("Schema").Where("verifier_did = ?", did).Find(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ProofRepository) CreateProofRequest(ctx context.Context, entity *proof.ProofRequest) (*proof.ProofRequest, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ProofRepository) UpdateProofRequest(ctx context.Context, entity *proof.ProofRequest, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProofRepository) FindProofResponseByPublicId(ctx context.Context, publicId string) (*proof.ProofResponse, error) {
	var entity proof.ProofResponse
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ProofRepository) FindAllProofResponses(ctx context.Context) ([]*proof.ProofResponse, error) {
	var entity []*proof.ProofResponse
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ProofRepository) CreateProofResponse(ctx context.Context, entity *proof.ProofResponse) (*proof.ProofResponse, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
