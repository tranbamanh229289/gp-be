package kafka

import (
	"be/config"
	"be/pkg/logger"
	"fmt"
	"sync"
)

type Manager struct {
	config    *config.Config
	logger    *logger.ZapLogger
	producers map[string]*Producer
	consumers map[string]*Consumer
	mu        sync.RWMutex
}

func NewManager(config *config.Config, logger *logger.ZapLogger) *Manager {
	return &Manager{
		config:    config,
		logger:    logger,
		producers: make(map[string]*Producer),
		consumers: make(map[string]*Consumer),
	}
}

func (m *Manager) CreateProducer(clientID string) (*Producer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.producers[clientID]; exists {
		return nil, fmt.Errorf("producer %s already exists", clientID)
	}

	producer, err := NewProducer(m.config, m.logger, clientID)
	if err != nil {
		return nil, err
	}

	m.producers[clientID] = producer

	return producer, nil
}

func (m *Manager) GetProducer(clientID string) (*Producer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.producers[clientID]; exists {
		return nil, fmt.Errorf("producer %s already exists", clientID)
	}
	return m.producers[clientID], nil
}

func (m *Manager) CreateConsumer(groupID string) (*Consumer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.consumers[groupID]; exists {
		return nil, fmt.Errorf("consumer %s already exists", groupID)
	}

	consumer, err := NewConsumer(m.config, m.logger, groupID)
	if err != nil {
		return nil, err
	}
	m.consumers[groupID] = consumer
	return consumer, nil
}

func (m *Manager) GetConsumer(groupID string) (*Consumer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.consumers[groupID]; exists {
		return nil, fmt.Errorf("consumer %s already exists", groupID)
	}
	return m.consumers[groupID], nil
}
