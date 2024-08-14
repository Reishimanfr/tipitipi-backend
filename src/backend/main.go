package main

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/routes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port = "8080"
)

func main() {
	router := gin.Default()
	err := core.InitDatabase()

	if err != nil {
		panic("Failed to initialize database: " + err.Error())
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

	fmt.Println("Server started on http://localhost:" + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully, waiting for current operations to complete
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	} else {
		fmt.Println("Server shutting down...")
	}
}
