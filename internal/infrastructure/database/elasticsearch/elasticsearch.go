package elasticsearch

import (
	"be/config"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchDB struct {
	Client *elasticsearch.Client
}

func NewDB(cfg *config.Config)(*ElasticsearchDB, error){
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
	
	return &ElasticsearchDB{Client: client}, nil
}

func (es *ElasticsearchDB) Close() error {
	return nil;
}