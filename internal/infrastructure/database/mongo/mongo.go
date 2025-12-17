package mongo

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
	logger *logger.ZapLogger
}

func NewDB(cfg *config.Config, logger *logger.ZapLogger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Mongo.Timeout*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(cfg.Mongo.URI).
			SetMaxPoolSize(cfg.Mongo.MaxPoolSize).
			SetMinPoolSize(cfg.Mongo.MinPoolSize))

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to mongo: %s", err))

		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error(fmt.Sprintf("Failed to ping mongo: %s", err))

		return nil, err
	}

	logger.Info("Successfully connected to Mongo",
		zap.String("addresses", cfg.Mongo.URI))

	return &MongoDB{
		client: client,
		db:     client.Database(cfg.Mongo.Database),
		logger: logger,
	}, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
