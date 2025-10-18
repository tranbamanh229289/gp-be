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
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},

			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
				tls.CurveP384,
			},

			SessionTicketsDisabled: false,

			NextProtos: []string{"h2", "http/1.1"},

			ClientAuth: tls.NoClientCert,
			
			Renegotiation: tls.RenegotiateNever,
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