package postgres

import (
	"be/config"
	"be/pkg/logger"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	sqlDB *sql.DB
	gormDB *gorm.DB
	logger *logger.ZapLogger
}

func NewDB(cfg *config.Config, logger *logger.ZapLogger) (*PostgresDB, error) {
	// native sql	
	dsn := config.GetPostgresDSN(cfg)

	sqlDB, err := sql.Open(cfg.Postgres.Driver, dsn)
	if err != nil {
		logger.Error("Failed to connect to database:", 
			zap.Error(err), 
			zap.String("addresses", dsn))
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxConnections)

	pingErr := sqlDB.Ping()
	if pingErr != nil {
		logger.Error("Failed to ping to database:", 
			zap.Error(pingErr), 
			zap.String("addresses", dsn))
		return nil, err
	}

	// gorm
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		logger.Error("Failed to connect to database:", 
			zap.Error(err), 
			zap.String("addresses", dsn))
		return nil, err
	}

	logger.Info("Successfully to connect to database:",  
			zap.String("addresses", dsn))
	
	return &PostgresDB{sqlDB: sqlDB, gormDB: gormDB, logger: logger}, nil
}

func (p *PostgresDB) GetGormDB() *gorm.DB {
	return p.gormDB
}

func (p *PostgresDB) Close() error {
	if err := p.sqlDB.Close(); err != nil {
		p.logger.Error("Failed to close to database:", 
			zap.Error(err),)
		return err
	}
	return nil
}