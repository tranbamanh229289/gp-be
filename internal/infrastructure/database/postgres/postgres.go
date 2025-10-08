package postgres

import (
	"database/sql"
	"graduate-project/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	SQLDB *sql.DB
	GormDB *gorm.DB
}

func NewDB(cfg *config.Config) (*PostgresDB, error) {
	// native sql
	dsn := config.GetPostgresDSN(cfg)

	sqlDB, err := sql.Open(cfg.Postgres.Driver, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err	)
	}

	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxConnections)

	pingErr := sqlDB.Ping()
	if pingErr != nil {
		log.Fatal("Failed to ping database:", pingErr)
	}

	// gorm
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err	)
	}
	
	return &PostgresDB{SQLDB: sqlDB, GormDB: gormDB}, nil
}

func (p *PostgresDB) Close() error {
	if err := p.SQLDB.Close(); err != nil {
		log.Fatal("Failed to close db")
	}
	return nil
}