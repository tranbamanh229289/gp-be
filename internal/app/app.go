package app

import (
	"be/config"
	"be/internal/infrastructure/cache/redis"
	"be/internal/infrastructure/database/elasticsearch"
	"be/internal/infrastructure/database/mongo"
	"be/internal/infrastructure/database/postgres"
	"be/internal/infrastructure/message_queue/rabbitmq"
	"be/internal/transport/http/middleware"
	"be/internal/transport/http/router"
	"be/pkg/fluent"
	"be/pkg/logger"
)

type App struct {
	Config *config.Config
	Router *router.Router
	Middleware *middleware.Middleware
	Server *Server
	Log *logger.ZapLogger
	Fluent *fluent.Fluent
	Postgres *postgres.PostgresDB
	Mongo *mongo.MongoDB
	Elasticsearch *elasticsearch.ElasticsearchDB
	Redis *redis.RedisCache
	RabbitMQQueue *rabbitmq.RabbitQueue
	RabbitMQConsumer *rabbitmq.Consumer
	RabbitMQProducer *rabbitmq.Producer
}
