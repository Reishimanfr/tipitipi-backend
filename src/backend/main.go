package main

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/flags"
	"bash06/strona-fundacja/src/backend/srv"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	log, err := core.InitLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", core.RandStr(128))
	}

	if !*flags.Dev {
		// Disable gin's debug logs
		gin.SetMode(gin.ReleaseMode)
	}

	serverConfig := &srv.ServerConfig{
		Port: *flags.Port,
		HttpConfig: &http.Server{
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
		CorsConfig: &cors.Config{
			AllowMethods:           []string{"HEAD", "POST", "DELETE", "PATCH", "GET"},
			AllowHeaders:           []string{"Content-Type", "Authorization"},
			AllowCredentials:       true,
			AllowFiles:             false,
			AllowAllOrigins:        true,
			AllowWebSockets:        false,
			AllowBrowserExtensions: false,
		},
	}

	server, err := srv.New(serverConfig)
	if err != nil {
		log.Fatal("Failed to initialize server", zap.Error(err))
	}

	server.InitHandler()

	go func() {
		if !*flags.Secure {
			if err := server.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("Failed to start server", zap.Error(err))
			}
		} else {
			if err := server.Http.ListenAndServeTLS(*flags.CertFilePath, *flags.KeyFilePath); err != nil && err != http.ErrServerClosed {
				log.Fatal("Failed to start secure server", zap.Error(err))
			}
		}

	}()

	log.Info("Server started. Listening to requests...")

	// Keeps the server from shutting itself down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Http.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Warn("Gracefully shutting down backend server...")
}
