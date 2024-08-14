package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.HEAD("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Create an http.Server with the Gin router
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server started on http://localhost:8080")

	// Wait for an interrupt signal to shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")

	// Create a deadline to wait for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully, waiting for current operations to complete
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	} else {
		fmt.Println("Server exiting")
	}
}
