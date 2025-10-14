package api

import (
	"be/config"
	"be/internal/app"
	"be/internal/transport/http/middleware"
	"be/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main(){
	// Initialize app
	app, err := app.InitializeApplication()

	if err != nil {
		log.Fatal("Failed to initialize application %v", err)
	}
	defer app.Log.Sync()

	engine := gin.New()

	// Setup global middleware
	SetupGlobalMiddlewares(app.Config, app.Log, engine)

	// Setup router
	SetupRoutes(engine)

	// Create HTTP Server 
	server := NewServer(app.Config, app.Log, engine)
	

	go func() {
		app.Log.Info("Starting server on: " + server.GetHttpServer().Addr)
		err := server.Run(app.Config)
		if err != nil && err != http.ErrServerClosed {
			app.Log.Info("Server closed!!! ")
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func SetupGlobalMiddlewares(cfg *config.Config, logger *logger.ZapLogger, engine *gin.Engine) {
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.ErrorHandlingMiddleware())
	engine.Use(middleware.LogMiddleware(logger))
}

func SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")
	{
		apiGroup.Use()
	}
}