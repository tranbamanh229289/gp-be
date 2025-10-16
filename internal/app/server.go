package app

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	logger *logger.ZapLogger
}

func NewServer(cfg *config.Config, logger *logger.ZapLogger) *Server{
	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Server.Host,cfg.Server.Port),
		Handler: nil,
		ReadTimeout: cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout: cfg.Server.Timeout,
	}

	if cfg.TLS.Enabled {
		httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
	}
	return &Server{httpServer: httpServer, logger: logger}
}

func (server *Server) SetHandler(engine *gin.Engine) {
	server.httpServer.Handler = engine
}

func(server *Server) Run(cfg *config.Config) error {
	if server.httpServer.TLSConfig != nil {
			server.logger.Info("TLS enabled")
			return server.httpServer.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
	}
	server.logger.Info("TLS disabled")
	return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}

func (server *Server) GetHttpServer() *http.Server{
	return server.httpServer
}