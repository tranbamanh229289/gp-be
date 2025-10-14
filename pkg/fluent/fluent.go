package fluent

import (
	"be/config"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type Fluent struct {
	logger *fluent.Fluent
	tag string
	mu sync.Mutex
}

func NewFluent(cfg *config.Config) (*Fluent, error){
	logger, err := fluent.New(fluent.Config{
		FluentHost: cfg.Fluent.Host,
		FluentPort: cfg.Fluent.Port,
		FluentNetwork: cfg.Fluent.Protocol,
		Timeout: cfg.Fluent.Timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create fluent logger %w", err)
	}

	return &Fluent{logger: logger, tag: cfg.App.Name}, nil
}

func (f *Fluent) Write(p []byte) (int, error){
	f.mu.Lock()
	defer f.mu.Unlock()
	logData := make(map[string]interface{})
	
	if err := parseLog(p, logData); err != nil {
		logData["raw"] = string(p)
	}

	err := f.logger.Post(f.tag, logData)
	return len(p), err
}	

func (f *Fluent) Sync() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.logger != nil {
		return f.logger.Close()
	}
	return nil
}

func (f *Fluent) Close() error {
	return f.Sync()
}

func (f *Fluent) GetFluent() *fluent.Fluent {
	return f.logger
}

func parseLog(data []byte, log map[string] interface{}) error {
	return json.Unmarshal(data, &log)
}

