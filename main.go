package main

import (
	"bash06/tipitipi-backend/core"
	"bash06/tipitipi-backend/flags"
	"bash06/tipitipi-backend/srv"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()

	log, err := core.InitLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	db, err := core.InitDb()
	if err != nil {
		log.Fatal("Failed to initialize database:", zap.Error(err))
	}

	if !*flags.Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	server, err := srv.New(&srv.ServerConfig{
		CorsConfig: &cors.Config{
			AllowMethods:           []string{"HEAD", "POST", "DELETE", "PATCH", "GET"},
			AllowHeaders:           []string{"Content-Type", "Authorization"},
			AllowOrigins:           []string{"http://localhost*", "https://tipitipi.pl"},
			AllowCredentials:       true,
			AllowFiles:             false,
			AllowWebSockets:        false,
			AllowBrowserExtensions: false,
			AllowWildcard:          true,
		},
		HttpConfig: &http.Server{},
		Logger:     log,
		Db:         db,
		Port:       flags.Port,
	})
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

	log.Info("Gracefully shutting down backend server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Http.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown:", zap.Error(err))
	}

	if <-ctx.Done(); true {
		log.Error("Timed out after 5 seconds")
	}

	log.Info("Server exiting")
}
