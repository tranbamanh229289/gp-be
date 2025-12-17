package elasticsearch

import (
	"be/config"
	"be/pkg/logger"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type ElasticsearchDB struct {
	client *elasticsearch.Client
	logger *logger.ZapLogger
}

func NewDB(cfg *config.Config, logger *logger.ZapLogger) (*ElasticsearchDB, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.GetElasticsearchAddress(cfg),
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create elasticsearch client: %s", zap.Error(err)))

		return nil, err
	}

	res, err := client.Info()

	if err != nil {
		logger.Error("Failed to get elasticsearch info: %s", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	logger.Info("Successfully connected to Elasticsearch",
		zap.Strings("addresses", config.GetElasticsearchAddress(cfg)))

	return &ElasticsearchDB{client: client, logger: logger}, nil
}

func (es *ElasticsearchDB) Close() error {
	return nil
}
