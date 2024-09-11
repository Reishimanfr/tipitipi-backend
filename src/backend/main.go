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
)

var (
	port   = "2333"
	server *http.Server
)

func startServer(router *gin.Engine) {
	server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Failed to start server: " + err.Error())
		}
	}()
}

func stopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic("Server forced to shutdown: " + err.Error())
	}
}

func setupRouter(db *core.Database, testing bool) *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = false

	if !testing {
		router.Use(middleware.RateLimiterMiddleware(middleware.NewRateLimiter(5, 10)))
		router.Use(middleware.FileSizeLimiterMiddleware(3000000000))

		router.Use(cors.New(cors.Config{
			AllowOrigins:           []string{"http://localhost:5173"},
			AllowMethods:           []string{"HEAD", "POST", "DELETE", "PATCH", "GET"},
			AllowHeaders:           []string{"Origin", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
			AllowCredentials:       true,
			AllowFiles:             false,
			AllowWebSockets:        false,
			AllowBrowserExtensions: false,
		}))
	}

	routes.NewHandler(&routes.Config{
		Router: router,
	}, db)

	return router
}

func main() {
	log, loggerErr := core.InitLogger()

	if loggerErr != nil {
		panic("Failed to initialize logger: " + loggerErr.Error())
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic("No JWT secret provided in .env file")
	}

	if os.Getenv("BACKEND_PORT") != "" {
		port = os.Getenv("BACKEND_PORT")
	}

	if os.Getenv("DEV") != "true" {
		// Disable gin's debug logs
		gin.SetMode(gin.ReleaseMode)
	}

	db := core.Database{
		Memory: false,
	}
	db.Init()

	r := setupRouter(&db, false)
	startServer(r)

	log.Info("Server started on http://localhost:" + port)

	// Keeps the server from shutting itself down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	stopServer()

	log.Warn("Gracefully shutting down backend server...")
}
