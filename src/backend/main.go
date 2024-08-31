package main

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/middleware"
	"bash06/strona-fundacja/src/backend/routes"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var port string

func main() {
	godotenv.Load(".env")

	log, loggerErr := core.InitLogger()

	if loggerErr != nil {
		panic("Failed to initialize logger: " + loggerErr.Error())
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic("No JWT secret provided in .env file")
	}

	port = os.Getenv("BACKEND_PORT")

	if port == "" {
		log.Warn("No backend port provided in .env, setting port to default value (2333)")
		port = "2333"
	}

	if os.Getenv("DEV") != "true" {
		// Disable gin's debug logs
		gin.SetMode(gin.ReleaseMode)
	}

	db := core.Database{}
	db.Init()

	router := gin.Default()

	router.Use(middleware.RateLimiterMiddleware(middleware.NewRateLimiter(5, 10)))

	// TODO: set this up correctly
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	routes.NewHandler(&routes.Config{
		Router: router,
	}, &db)

	// TODO: set this correctly also
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
