package postgres

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	sqlDB   *sql.DB
	gormDB  *gorm.DB
	pgxPool *pgxpool.Pool
	logger  *logger.ZapLogger
}

func NewDB(cfg *config.Config, logger *logger.ZapLogger) (*PostgresDB, error) {
	dsn := cfg.GetPostgresDSN()
	ctx := context.Background()
	// pgx
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("failed to create pgx pool")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to ping postgres")
		return nil, err
	}

	// native sql
	sqlDB := stdlib.OpenDBFromPool(pool)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)

	// gorm
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %s", err))
		return nil, err
	}

	logger.Info("Successfully to connect to database:",
		zap.String("dsn", dsn))

	return &PostgresDB{sqlDB: sqlDB, gormDB: gormDB, pgxPool: pool, logger: logger}, nil
}

func (p *PostgresDB) Close() error {
	if p.sqlDB != nil {
		err := p.sqlDB.Close()
		if err != nil {
			p.logger.Error(fmt.Sprintf("Failed to close to database: %s", err))
			return err
		}
	}
	if p.pgxPool != nil {
		p.pgxPool.Close()
	}
	return nil
}

func (p *PostgresDB) GetGormDB() *gorm.DB {
	return p.gormDB
}

func (p *PostgresDB) GetPgxPool() *pgxpool.Pool {
	return p.pgxPool
}
