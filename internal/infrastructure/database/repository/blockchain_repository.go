package repository

import (
	"be/internal/domain/blockchain"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type BlockchainRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewBlockchainRepository(logger *logger.ZapLogger, db *postgres.PostgresDB) blockchain.IBlockchainRepository {
	return &BlockchainRepository{
		logger: logger,
		db:     db,
	}
}

func (r *BlockchainRepository) GetLastBlock(ctx context.Context) uint64 {
	return 10
}

func (r *BlockchainRepository) SaveBlocks(ctx context.Context, blocks []*blockchain.Block) error {
	return nil
}

func (r *BlockchainRepository) SaveLastBlockNumber(ctx context.Context, blockNumber uint64) error {
	return nil
}

func (r *BlockchainRepository) SaveTransactions(ctx context.Context, txs []*blockchain.Transaction) error {
	return nil
}

func (r *BlockchainRepository) SaveEvents(ctx context.Context, events []*blockchain.Event) error {
	return nil
}
