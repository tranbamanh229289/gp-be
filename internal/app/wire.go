//go:build wireinject
// +build wireinject

package app

import (
	"be/config"
	"be/internal/infrastructure/blockchain/ether"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/elasticsearch"
	"be/internal/infrastructure/database/mongo"
	"be/internal/infrastructure/database/postgres"
	"be/internal/infrastructure/database/repository"
	"be/internal/infrastructure/message_queue/rabbitmq"
	"be/internal/service"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"
	"be/internal/transport/http/router"
	"be/pkg/fluent"
	"be/pkg/logger"

	"github.com/google/wire"
)

// Log Set
var logSet = wire.NewSet(fluent.NewFluent, logger.NewLogger)

// Config Set
var configSet = wire.NewSet(config.NewConfig)

// Infra Set
var dbSet = wire.NewSet(postgres.NewDB, mongo.NewDB, elasticsearch.NewDB)
var cacheSet = wire.NewSet(redis.NewCache)
var queueSet = wire.NewSet(rabbitmq.NewQueue, rabbitmq.NewConsumer, rabbitmq.NewProducer)
var etherSet = wire.NewSet(ether.NewEther)

// Handler Set
var handlerSet = wire.NewSet(handler.NewAuthHandler)

// Service Set
var serviceSet = wire.NewSet(service.NewAuthService, service.BlockchainService)

// Repository Set
var repositorySet = wire.NewSet(repository.NewUserRepository)

// Router Set
var routerSet = wire.NewSet(router.NewRouter)

// Middleware Set
var middlewareSet = wire.NewSet(middleware.NewMiddleware)

// Server Set
var serverSet = wire.NewSet(NewServer)

func InitializeApplication() (App, error) {
    panic(wire.Build(
        configSet,
        logSet,
        dbSet,
        cacheSet,
        queueSet,
        handlerSet,
        serviceSet,
        repositorySet,
				routerSet,
				middlewareSet,
				serverSet,
        wire.Struct(new(App), "*"),
    ))
}