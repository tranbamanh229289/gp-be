package app

import (
	"be/config"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/mongo"
	"be/internal/infrastructure/database/postgres"
	"be/internal/infrastructure/message_queue/rabbitmq"
	"be/pkg/fluent"
	"be/pkg/logger"
)

type App struct {
	Config *config.Config
	Log *logger.ZapLogger
	Fluent *fluent.Fluent
	Postgres *postgres.PostgresDB
	Mongo *mongo.MongoDB
	Redis *redis.RedisCache
	RabbitMQQueue *rabbitmq.RabbitQueue
	RabbitMQConsumer *rabbitmq.Consumer
	RabbitMQProducer *rabbitmq.Producer
}
