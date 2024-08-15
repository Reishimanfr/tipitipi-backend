package main

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/routes"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	port = "8080"
)

func main() {
	router := gin.Default()
	dbError := core.InitDatabase()

	loggerErr := core.InitLogger()
	log := core.GetLogger()

	if dbError != nil {
		panic("Failed to initialize database: " + dbError.Error())
	}

	if loggerErr != nil {
		panic("Failed to initialize logger: " + loggerErr.Error())
	}

	// TODO: set this to prod if env variable says we're in docker
	// gin.SetMode(gin.ReleaseMode)

	routes.NewHandler(&routes.Config{
		Router: router,
	})

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     router,
		IdleTimeout: -1,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Failed to start server: " + err.Error())
		}
	}()

	log.Info("Server started on http://localhost:" + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Warn("Gracefully shutting down backend server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}
}
