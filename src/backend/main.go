package main

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/srv"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	port = "2333"
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

	serverConfig := &srv.ServerConfig{
		Port:    port,
		Testing: false,
		HttpConfig: &http.Server{
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
		AwsConfig: &aws.Config{
			Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Region:           aws.String(region),
			Endpoint:         aws.String(endpoint),
			S3ForcePathStyle: aws.Bool(true),
		},
		CorsConfig: &cors.Config{
			AllowMethods:           []string{"HEAD", "POST", "DELETE", "PATCH", "GET"},
			AllowHeaders:           []string{"*"},
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
		if err := server.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	log.Info("Server started on http://localhost:" + port)

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
