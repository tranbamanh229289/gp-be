package repository

import (
	"be/config"
	"be/internal/infrastructure/database/postgres"
	"context"

	"github.com/iden3/go-merkletree-sql/v2"
)

type IMTRepository interface {
	NewMerkleTree(ctx context.Context) (*merkletree.MerkleTree, uint64, error)
	LoadMerkleTree(ctx context.Context, mtID uint64) (*merkletree.MerkleTree, error)
}

type MTRepository struct {
	config *config.Config
	db     *postgres.PostgresDB
}

func NewMerkletreeRepository(config *config.Config, db *postgres.PostgresDB) IMTRepository {
	return &MTRepository{config: config, db: db}
}

func (r *MTRepository) NewMerkleTree(ctx context.Context) (*merkletree.MerkleTree, uint64, error) {
	mtId, err := NextMTID(ctx, r.db.GetPgxPool())

	storage := NewSqlStorage(r.db.GetPgxPool(), mtId)
	mt, err := merkletree.NewMerkleTree(ctx, storage, r.config.Circuit.MTLevel)

	return mt, mtId, err
}

func (r *MTRepository) LoadMerkleTree(ctx context.Context, mtID uint64) (*merkletree.MerkleTree, error) {
	storage := NewSqlStorage(r.db.GetPgxPool(), mtID)
	mt, err := merkletree.NewMerkleTree(ctx, storage, r.config.Circuit.MTLevel)

	return mt, err
}
