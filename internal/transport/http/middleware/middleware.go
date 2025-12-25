package middleware

import (
	"be/config"
	"be/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	config *config.Config
	logger *logger.ZapLogger
}

func NewMiddleware(cfg *config.Config, logger *logger.ZapLogger) *Middleware {
	return &Middleware{config: cfg, logger: logger}
}

func (m *Middleware) SetupGlobalMiddlewares(engine *gin.Engine) {
	m.CORSMiddleware(engine)
	m.LogMiddleware(engine)
}
