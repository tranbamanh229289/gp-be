package service

import "be/config"

type ISchemaService interface {
}
type SchemaService struct {
	config *config.Config
}

func NewSchemaService(config *config.Config) ISchemaService {
	return &SchemaService{
		config: config,
	}
}
