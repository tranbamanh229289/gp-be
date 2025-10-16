package repository

import (
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
)

type BlockchainRepository struct{
	db *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewBlockchainRepository(logger *logger.ZapLogger, db *postgres.PostgresDB) *BlockchainRepository {
	return &BlockchainRepository{
		logger: logger,
		db: db,
	}
}
