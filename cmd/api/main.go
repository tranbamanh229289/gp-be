package main

import (
	"be/internal/app"
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
		log.Printf("Failed to initialize application %s", err)
	}
	defer app.Log.Sync()

	engine := gin.New()

	// Setup global middleware
	app.Middleware.SetupGlobalMiddlewares(engine)

	// Setup router
	app.Router.SetupRoutes(engine)

	// Create HTTP Server 
	app.Server.SetHandler(engine)
	

	go func() {
		app.Log.Info("Starting server on: " + app.Server.GetHttpServer().Addr)
		err := app.Server.Run(app.Config)
		if err != nil && err != http.ErrServerClosed {
			app.Log.Info("Server closed!!! ")
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	app.Server.Shutdown(ctx)
}

