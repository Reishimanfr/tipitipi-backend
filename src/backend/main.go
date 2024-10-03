package main

import (
	ovh "bash06/strona-fundacja/src/backend/aws"
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

const (
	err_jwt_not_set                  = "JWT secret not set"
	err_aws_access_missing           = "AWS access token missing"
	err_aws_secret_missing           = "AWS secret token missing"
	err_aws_region_missing           = "AWS region missing (waw?)"
	err_aws_endpoint_missing         = "AWS endpoint missing"
	err_aws_worker_init_failure      = "AWS worker failed to initialize"
	err_aws_blog_bucket_name_missing = "Blog bucket name missing"
)

func main() {
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	blogBucket := os.Getenv("AWS_BLOG_BUCKET_NAME")

	if accessKey == "" {
		panic(err_aws_access_missing)
	}

	if secretKey == "" {
		panic(err_aws_secret_missing)
	}

	if endpoint == "" {
		panic(err_aws_endpoint_missing)
	}

	if region == "" {
		panic(err_aws_region_missing)
	}

	if blogBucket == "" {
		panic(err_aws_blog_bucket_name_missing)
	}

	log, loggerErr := core.InitLogger()

	if loggerErr != nil {
		panic("Failed to initialize logger: " + loggerErr.Error())
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic(err_jwt_not_set)
	}

	if os.Getenv("BACKEND_PORT") != "" {
		port = os.Getenv("BACKEND_PORT")
	}

	if os.Getenv("DEV") != "true" {
		// Disable gin's debug logs
		gin.SetMode(gin.ReleaseMode)
	}

	db := core.Database{}
	db.Init()

	worker, err := ovh.NewWorker(accessKey, secretKey, region, endpoint)
	if err != nil {
		panic(err_aws_worker_init_failure)
	}

	router := gin.Default()
	router.RedirectTrailingSlash = false

	router.Use(middleware.RateLimiterMiddleware(middleware.NewRateLimiter(5, 10)))
	router.Use(middleware.FileSizeLimiterMiddleware(3000000000)) // TODO: fix this

	router.Use(cors.New(cors.Config{
		AllowOrigins:           []string{"http://localhost:5173"},
		AllowMethods:           []string{"HEAD", "POST", "DELETE", "PATCH", "GET"},
		AllowHeaders:           []string{"Origin", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
		AllowCredentials:       true,
		AllowFiles:             false,
		AllowWebSockets:        false,
		AllowBrowserExtensions: false,
	}))

	routes.NewHandler(&routes.Config{
		Router: router,
	}, &db, worker)

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

	log.Info("Server started on http://localhost:" + port)

	// Keeps the server from shutting itself down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic("Server forced to shutdown: " + err.Error())
	}

	log.Warn("Gracefully shutting down backend server...")
}
