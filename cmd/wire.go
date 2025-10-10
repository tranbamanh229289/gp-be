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

var configSet = wire.NewSet(config.NewConfig)

var dbSet = wire.NewSet(postgres.NewDB, mongo.NewDB, elasticsearch.NewDB)
var cacheSet = wire.NewSet(redis.NewCache)
var queueSet = wire.NewSet(rabbitmq.NewQueue, rabbitmq.NewConsumer, rabbitmq.NewProducer)

var domainSet = wire.NewSet()

var serviceSet = wire.NewSet()

var pkgSet = wire.NewSet()

func InitializeDependency() error {
	wire.Build(configSet, dbSet, queueSet, domainSet, serviceSet, pkgSet)
	return nil
}