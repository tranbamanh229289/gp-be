package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name           string
	Version        string
	Env            string
	AllowedOrigins string
}
type ServerConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type JwtConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
type ZapConfig struct {
	Level    string
	Encoding string
}

type PostgresConfig struct {
	Driver             string
	Host               string
	Port               int
	User               string
	Password           string
	Name               string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
}

type MongoConfig struct {
	URI         string
	Username    string
	Password    string
	Database    string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout     time.Duration
}

type ElasticsearchConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type RedisConfig struct {
	Host           string
	Port           int
	Password       string
	DB             int
	MaxConnections int
	Timeout        time.Duration
}

type RabbitMQConfig struct {
	Host              string
	Port              int
	UserName          string
	Password          string
	VHost             string
	PrefetchCount     int
	ConnectionTimeout int
	HeartBeat         int
	ReconnectDelay    int
}

type KafkaProducerConfig struct {
	BootstrapServers  string
	Acks              string
	CompressionType   string
	LingerMs          int
	BatchSize         int
	EnableIdempotence bool
	Retries           int
	// ClientID          string
}

type KafkaConsumerConfig struct {
	BootstrapServers     string
	AutoOffsetReset      string
	SessionTimeoutMs     int
	EnableAutoCommit     bool
	AutoCommitIntervalMs int
	// GroupID              string
	// Topics               []string
}

type KafkaConfig struct {
	Producer KafkaProducerConfig
	Consumer KafkaConsumerConfig
}

type FluentConfig struct {
	Host     string
	Port     int
	Protocol string // tcp, udp
	Timeout  time.Duration
}

type BlockchainConfig struct {
	RPC           string
	Resolver      string
	StateContract string
}

type PinataConfig struct {
	APIKey       string
	APISecretKey string
	JWTKey       string
	Endpoint     string
	GatewayURL   string
}
type CircuitConfig struct {
	MTLevel        int
	MTLevelOnChain int
	MTLevelClaim   int

	VerifyingKey string
}

type Iden3Config struct {
	VerifierPrivateKey string
}
type CronConfig struct {
}

type Config struct {
	App           AppConfig
	Server        ServerConfig
	TLS           TLSConfig
	Postgres      PostgresConfig
	Mongo         MongoConfig
	Elasticsearch ElasticsearchConfig
	Redis         RedisConfig
	RabbitMQ      RabbitMQConfig
	Kafka         KafkaConfig
	JWT           JwtConfig
	Zap           ZapConfig
	Fluent        FluentConfig
	Cron          CronConfig
	Blockchain    BlockchainConfig
	IPFS          PinataConfig
	Circuit       CircuitConfig
	Iden3         Iden3Config
}

func NewConfig() (*Config, error) {
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
		App: AppConfig{
			Name:           viper.GetString("app.name"),
			Version:        viper.GetString("app.version"),
			Env:            viper.GetString("app.env"),
			AllowedOrigins: viper.GetString("app.allowed_origins"),
		},
		Server: ServerConfig{
			Host:    viper.GetString("server.host"),
			Port:    viper.GetInt("server.port"),
			Timeout: viper.GetDuration("server.timeout"),
		},
		TLS: TLSConfig{
			Enabled:  viper.GetBool("tls.enabled"),
			CertFile: viper.GetString("tls.cert_file"),
			KeyFile:  viper.GetString("tls.key_file"),
		},
		Postgres: PostgresConfig{
			Host:               viper.GetString("postgres.host"),
			Port:               viper.GetInt("postgres.port"),
			User:               viper.GetString("postgres.user"),
			Password:           viper.GetString("postgres.password"),
			Name:               viper.GetString("postgres.name"),
			Driver:             viper.GetString("postgres.driver"),
			SSLMode:            viper.GetString("postgres.sslmode"),
			MaxConnections:     viper.GetInt("postgres.max_connections"),
			MaxIdleConnections: viper.GetInt("postgres.max_idle_connections"),
		},
		Mongo: MongoConfig{
			URI:         viper.GetString("mongo.uri"),
			Database:    viper.GetString("mongo.database"),
			Username:    viper.GetString("mongo.username"),
			Password:    viper.GetString("mongo.password"),
			MaxPoolSize: viper.GetUint64("mongo.max_pool_size"),
			MinPoolSize: viper.GetUint64("mongo.min_pool_size"),
			Timeout:     viper.GetDuration("mongo.timeout"),
		},
		Elasticsearch: ElasticsearchConfig{
			Host:     viper.GetString("elasticsearch.host"),
			Port:     viper.GetInt("elasticsearch.port"),
			Username: viper.GetString("elasticsearch.username"),
			Password: viper.GetString("elasticsearch.password"),
		},
		Redis: RedisConfig{
			Host:           viper.GetString("redis.host"),
			Port:           viper.GetInt("redis.port"),
			Password:       viper.GetString("redis.password"),
			DB:             viper.GetInt("redis.db"),
			MaxConnections: viper.GetInt("redis.max_connections"),
			Timeout:        viper.GetDuration("redis.timeout"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:              viper.GetString("rabbitmq.host"),
			Port:              viper.GetInt("rabbitmq.port"),
			UserName:          viper.GetString("rabbitmq.username"),
			Password:          viper.GetString("rabbitmq.password"),
			VHost:             viper.GetString("rabbitmq.vhost"),
			PrefetchCount:     viper.GetInt("rabbitmq.prefetch_count"),
			ConnectionTimeout: viper.GetInt("rabbitmq.connection_timeout"),
			HeartBeat:         viper.GetInt("rabbitmq.heartbeat"),
			ReconnectDelay:    viper.GetInt("rabbitmq.reconnect_delay"),
		},
		Kafka: KafkaConfig{
			Producer: KafkaProducerConfig{
				BootstrapServers: viper.GetString("kafka.producer.bootstrap_servers"),
				// ClientID:          viper.GetString("kafka.producer.client_id"),
				Acks:              viper.GetString("kafka.producer.acks"),
				CompressionType:   viper.GetString("kafka.producer.compress_type"),
				LingerMs:          viper.GetInt("kafka.producer.linger_ms"),
				BatchSize:         viper.GetInt("kafka.producer.batch_size"),
				EnableIdempotence: viper.GetBool("kafka.producer.enable_idempotence"),
				Retries:           viper.GetInt("kafka.producer.retries"),
			},
			Consumer: KafkaConsumerConfig{
				BootstrapServers: viper.GetString("kafka.producer.bootstrap_servers"),
				// GroupID:              viper.GetString("kafka.producer.group_id"),
				// Topics:               viper.GetStringSlice("kafka.producer.topics"),
				AutoOffsetReset:      viper.GetString("kafka.producer.auto_offset_reset"),
				SessionTimeoutMs:     viper.GetInt("kafka.producer.session_timeout_ms"),
				EnableAutoCommit:     viper.GetBool("kafka.producer.enable_auto_commit"),
				AutoCommitIntervalMs: viper.GetInt("kafka.producer.auto_commit_interval_ms"),
			},
		},
		JWT: JwtConfig{
			Secret:          viper.GetString("jwt.secret"),
			AccessTokenTTL:  viper.GetDuration("jwt.access_token_ttl"),
			RefreshTokenTTL: viper.GetDuration("jwt.refresh_token_ttl"),
		},
		Zap: ZapConfig{
			Level:    viper.GetString("zap.level"),
			Encoding: viper.GetString("zap.encoding"),
		},
		Fluent: FluentConfig{
			Host:     viper.GetString("fluent.host"),
			Port:     viper.GetInt("fluent.port"),
			Protocol: viper.GetString("fluent.protocol"),
			Timeout:  viper.GetDuration("fluent.timeout"),
		},
		Blockchain: BlockchainConfig{
			RPC:           viper.GetString("blockchain.polygon.amoy.rpc"),
			Resolver:      viper.GetString("blockchain.polygon.amoy.resolver"),
			StateContract: viper.GetString("blockchain.polygon.amoy.contract.state"),
		},
		IPFS: PinataConfig{
			APIKey:       viper.GetString("ipfs.pinata.api_key"),
			APISecretKey: viper.GetString("ipfs.pinata.api_secret_key"),
			JWTKey:       viper.GetString("ipfs.pinata.jwt_key"),
			Endpoint:     viper.GetString("ipfs.pinata.endpoint"),
			GatewayURL:   viper.GetString("ipfs.pinata.gateway_url"),
		},
		Circuit: CircuitConfig{
			MTLevel:        viper.GetInt("circuit.mt_level"),
			MTLevelOnChain: viper.GetInt("circuit.mt_level_onchain"),
			MTLevelClaim:   viper.GetInt("circuit.mt_level_claim"),
			VerifyingKey:   viper.GetString("circuit.verifying_key"),
		},
		Iden3: Iden3Config{
			VerifierPrivateKey: viper.GetString("iden3.verifier.private_key"),
		},
	}
	return config, nil
}

func (config *Config) GetBaseURL() string {
	return fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
}

func (config *Config) GetPostgresDSN() string {
	host := config.Postgres.Host
	port := config.Postgres.Port
	user := config.Postgres.User
	password := config.Postgres.Password
	dbName := config.Postgres.Name
	sslmode := config.Postgres.SSLMode

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, strconv.Itoa(port), user, password, dbName, sslmode)
}

func (config *Config) GetRabbitMQDSN() string {
	host := config.RabbitMQ.Host
	port := config.RabbitMQ.Port
	username := config.RabbitMQ.UserName
	password := config.RabbitMQ.Password
	vhost := config.RabbitMQ.VHost

	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
}

func (config *Config) GetElasticsearchAddress() []string {
	return []string{fmt.Sprintf("http://%s:%s", config.Elasticsearch.Host, strconv.Itoa(config.Elasticsearch.Port))}
}
