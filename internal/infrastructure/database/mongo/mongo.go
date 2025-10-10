package mongo

import (
	"be/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB *mongo.Database
}

func NewDB(cfg *config.Config) (*MongoDB, error){
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Mongo.Timeout * time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx, 
		options.Client().
		ApplyURI(cfg.Mongo.URI).
		SetMaxPoolSize(cfg.Mongo.MaxPoolSize).
		SetMinPoolSize(cfg.Mongo.MinPoolSize))

	if err != nil {
		log.Fatal("Failed to connect to mongo:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping mongo:", err)
	}

	return &MongoDB{
		Client: client,
		DB: client.Database(cfg.Mongo.Database),
	}, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}