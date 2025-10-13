package api

import (
	"be/config"
	"be/internal/transport/http/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)


func SetupGlobalMiddlewares(engine *gin.Engine) {
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.ErrorHandlingMiddleware())
}

func SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")
	{
		apiGroup.Use()
	}
}

func main(){
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config:%v", err)
	}

	engine := gin.New()

	// Setup global middleware
	SetupGlobalMiddlewares(engine)

	// Setup router
	SetupRoutes(engine)

	// Create HTTP Server 
	server := NewServer(cfg, engine)
	
	go func() {
		err := server.Run(cfg)
		if err != nil && err != http.ErrServerClosed {
			log.Println("Server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
