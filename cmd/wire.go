//go:build wireinject
// +build wireinject

package cmd

import (
	"be/config"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/elasticsearch"
	"be/internal/infrastructure/database/mongo"
	"be/internal/infrastructure/database/postgres"
	"be/internal/infrastructure/message_queue/rabbitmq"

	"github.com/google/wire"
)

// Infra Set
var dbSet = wire.NewSet(postgres.NewDB, mongo.NewDB, elasticsearch.NewDB)
var cacheSet = wire.NewSet(redis.NewCache)
var queueSet = wire.NewSet(rabbitmq.NewQueue, rabbitmq.NewConsumer, rabbitmq.NewProducer)

// Domain Set
var domainSet = wire.NewSet()

// Service Set
var serviceSet = wire.NewSet()

// Handler Set

func InitializeDependency(config *config.Config) error {
	wire.Build(dbSet, queueSet, domainSet, serviceSet)
	return nil
}