package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port int
	Timeout time.Duration
}

type TLSConfig struct {
	Enabled bool
	CertFile string
	KeyFile string
}

type PostgresConfig struct {
	Driver string
	Host string
	Port int
	User string
	Password string
	Name string
	SSLMode string
	MaxConnections int
}

type MongoConfig struct {
	URI string
	Username string
	Password string
	Database string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout time.Duration
}

type ElasticsearchConfig struct {
	Host string
	Port int
	Username string
	Password string
}

type RedisConfig struct {
	Host string
	Port int
	Password string
	DB int
	MaxConnections int
	Timeout time.Duration
}

type RabbitMQConfig struct {
		Host string
		Port int
		UserName string
		Password string
		VHost string
		PrefetchCount int
		ConnectionTimeout int
		HeartBeat int
		ReconnectDelay int

}

type JwtConfig struct {
	Secret string
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
}

type Config struct {
	Server ServerConfig
	TLS TLSConfig
	Postgres PostgresConfig
	Mongo MongoConfig
	Elasticsearch ElasticsearchConfig
	Redis RedisConfig
	RabbitMQ RabbitMQConfig
	JwtConfig JwtConfig
}


func LoadConfig()(*Config, error) {
	// read .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// replace . to _ for env variables 
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// set config
	viper.SetConfigName("config_dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file:", err)
	}
	
	config := &Config{
		Server: ServerConfig{
			Host: viper.GetString("server.host"),
			Port: viper.GetInt("server.port"),
			Timeout: viper.GetDuration("server.timeout"),
		},
		TLS: TLSConfig{
			Enabled: viper.GetBool("tls.enabled"),
			CertFile: viper.GetString("tls.cert_file"),
			KeyFile: viper.GetString("tls.key_file"),
		},
		Postgres: PostgresConfig{
			Host: viper.GetString("postgres.host"),
			Port: viper.GetInt("postgres.port"),
			User: viper.GetString("postgres.user"),
			Password: viper.GetString("postgres.password"),
			Name: viper.GetString("postgres.name"),
			SSLMode: viper.GetString("postgres.sslmode"),
			MaxConnections: viper.GetInt("postgres.max_connections"),
		},
		Mongo: MongoConfig{
			URI: viper.GetString("mongo.uri"),
			Database: viper.GetString("mongo.database"),
			Username: viper.GetString("mongo.username"),
			Password: viper.GetString("mongo.password"),
			MaxPoolSize: viper.GetUint64("mongo.max_pool_size"),
			MinPoolSize: viper.GetUint64("mongo.min_pool_size"),
			Timeout: viper.GetDuration("mongo.timeout"),
		},
		Elasticsearch: ElasticsearchConfig{
			Host: viper.GetString("elasticsearch.host"),
			Port: viper.GetInt("elasticsearch.port"),
			Username: viper.GetString("elasticsearch.username"),
			Password: viper.GetString("elasticsearch.password"),
		},
		Redis: RedisConfig{
			Host: viper.GetString("redis.host"),
			Port: viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB: viper.GetInt("redis.db"),
			MaxConnections: viper.GetInt("redis.max_connections"),
			Timeout: viper.GetDuration("redis.timeout"),
		},
		RabbitMQ: RabbitMQConfig{
			Host: viper.GetString("rabbitmq.host"),
			Port: viper.GetInt("rabbitmq.port"),
			UserName: viper.GetString("rabbitmq.username"),
			Password: viper.GetString("rabbitmq.password"),
			VHost: viper.GetString("rabbitmq.vhost"),
			PrefetchCount: viper.GetInt("rabbitmq.prefetch_count"),
			ConnectionTimeout: viper.GetInt("rabbitmq.connection_timeout"),
			HeartBeat: viper.GetInt("rabbitmq.heartbeat"),
			ReconnectDelay: viper.GetInt("rabbitmq.reconnect_delay"),
		},
		JwtConfig: JwtConfig{
			Secret: viper.GetString("jwt.secret"),
			AccessTokenTTL: viper.GetDuration("jwt.access_token_ttl"),
			RefreshTokenTTL: viper.GetDuration("jwt.refresh_token_ttl"),
		},
		
	}
	return config, nil
}

func GetServerAddress(config *Config) string {
	return fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
}

func GetPostgresDSN(config *Config) string {
	host := config.Postgres.Host
	port := config.Postgres.Port
	user := config.Postgres.User
	password := config.Postgres.Password
	dbName := config.Postgres.Name
	sslmode := config.Postgres.SSLMode

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslmode)
}

func GetRabbitMQDSN(config *Config) string {
	host := config.RabbitMQ.Host
	port := config.RabbitMQ.Port
	username := config.RabbitMQ.UserName
	password := config.RabbitMQ.Password
	vhost := config.RabbitMQ.VHost

	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
}

func GetElasticsearchAddress(config *Config) []string {
	return []string{fmt.Sprintf("http://%s:%d", config.Elasticsearch.Host, config.Elasticsearch.Port)}
}

