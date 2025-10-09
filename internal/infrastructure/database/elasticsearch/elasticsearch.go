package elasticsearch

import (
	"be/config"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	Client *elasticsearch.Client
}

func NewElasticsearchClient(cfg *config.Config)(*ElasticsearchClient, error){
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.GetElasticsearchAddress(cfg),
		Username: cfg.Elasticsearch.Username,
		Password: cfg.Elasticsearch.Password,
	})

	if err != nil {
		log.Fatal("Failed to create elasticsearch client:", err)
	}

	res, err := client.Info();

	if  err != nil {
		log.Fatal("Failed to get elasticsearch info:", err)
	}
	defer res.Body.Close()
	log.Println("Elasticsearch info:", res)
	
	return &ElasticsearchClient{Client: client}, nil
}

func (ec *ElasticsearchClient) Close() error {
	return nil;
}