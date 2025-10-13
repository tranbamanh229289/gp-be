package api

import (
	"be/config"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server{
	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Server.Host,cfg.Server.Port),
		Handler: handler,
		ReadTimeout: cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout: cfg.Server.Timeout,
	}

	if cfg.TLS.Enabled {
		httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
	}
	return &Server{httpServer: httpServer}
}

func(server *Server) Run(cfg *config.Config) error {
	log.Printf("Starting server on %s", server.httpServer.Addr)
    if server.httpServer.TLSConfig != nil {
        log.Println("TLS enabled")
        return server.httpServer.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
    }
    log.Println("TLS disabled")
    return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}