//go:build wireinject
// +build wireinject

package app

import (
	"be/config"
	"be/internal/infrastructure/blockchain/ether"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/postgres"
	"be/internal/infrastructure/database/repository"
	"be/internal/service"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"
	"be/internal/transport/http/router"
	"be/pkg/logger"

	"github.com/google/wire"
)

// Log Set
var logSet = wire.NewSet(logger.NewLogger)

// Config Set
var configSet = wire.NewSet(config.NewConfig)

// Infra Set
var dbSet = wire.NewSet(postgres.NewDB)
var migrateSet = wire.NewSet()
var cacheSet = wire.NewSet(redis.NewCache)

// var queueSet = wire.NewSet(rabbitmq.NewQueue, rabbitmq.NewConsumer, rabbitmq.NewProducer)
var etherSet = wire.NewSet(ether.NewEther)

// Handler Set
var handlerSet = wire.NewSet(
	handler.NewAuthJWTHandler,
	handler.NewAuthZkHandler,
	handler.NewIssuerHandler,
	handler.NewHolderHandler,
	handler.NewVerifierHandler,
	handler.NewCredentialHandler,
)

// Service Set
var serviceSet = wire.NewSet(
	service.NewAuthJWTService,
	service.NewAuthZkService,
	service.NewIssuerService,
	service.NewHolderService,
	service.NewVerifierService,
	service.NewCredentialService,
)

// Repository Set
var repositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewBlockchainRepository,
	repository.NewCitizenIdentityRepository,
	repository.NewAcademicDegreeRepository,
	repository.NewDriverLicenseRepository,
	repository.NewHealthInsuranceRepository,
	repository.NewPassportRepository,
)

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
		// queueSet,
		repositorySet,
		serviceSet,
		handlerSet,
		routerSet,
		middlewareSet,
		serverSet,
		wire.Struct(new(App), "*"),
	))
}
