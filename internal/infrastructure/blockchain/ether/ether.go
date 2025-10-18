package ether

import (
	"be/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Ether struct {
	client *ethclient.Client
}

func NewEther(cfg *config.Config) (*Ether, error) {
	client, err := ethclient.Dial(cfg.Ether.RPC)
	if err != nil {
		return nil, err
	}
	return &Ether{client: client}, nil
}

func (e *Ether) GetClient() *ethclient.Client {
	return e.client
}