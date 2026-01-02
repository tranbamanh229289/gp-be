//go:build wireinject
// +build wireinject

package app

import (
	"be/config"
	"be/internal/infrastructure/blockchain/ether"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/postgres"
	"be/internal/infrastructure/database/repository"
	"be/internal/infrastructure/ipfs"
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
var ipfsSet = wire.NewSet(ipfs.NewPinata)

// var queueSet = wire.NewSet(rabbitmq.NewQueue, rabbitmq.NewConsumer, rabbitmq.NewProducer)
var etherSet = wire.NewSet(ether.NewEther)

// Handler Set
var handlerSet = wire.NewSet(
	handler.NewAuthJWTHandler,
	handler.NewAuthZkHandler,
	handler.NewDocumentHandler,
	handler.NewSchemaHandler,
	handler.NewCredentialHandler,
	handler.NewProofHandler,
)

// Service Set
var serviceSet = wire.NewSet(
	service.NewAuthJWTService,
	service.NewAuthZkService,
	service.NewCredentialService,
	service.NewDocumentService,
	service.NewProofService,
	service.NewSchemaService,
	service.NewIdentityService,
	service.NewVerifierService,
)

// Repository Set
var repositorySet = wire.NewSet(
	repository.NewAcademicDegreeRepository,
	repository.NewCitizenIdentityRepository,
	repository.NewCredentialRequestRepository,
	repository.NewDriverLicenseRepository,
	repository.NewHealthInsuranceRepository,
	repository.NewIdentityRepository,
	repository.NewMerkletreeRepository,
	repository.NewPassportRepository,
	repository.NewProofRepository,
	repository.NewSchemaAttributeRepository,
	repository.NewSchemaRepository,
	repository.NewStateTransitionRepository,
	repository.NewUserRepository,
	repository.NewVerifiableCredentialRepository,
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
		ipfsSet,
		repositorySet,
		serviceSet,
		handlerSet,
		routerSet,
		middlewareSet,
		serverSet,
		wire.Struct(new(App), "*"),
	))
}
