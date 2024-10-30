package main

import (
	"bash06/strona-fundacja/src/backend/srv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	server *srv.Server
	err    error
)

func setup() {
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

	if os.Getenv("JWT_SECRET") == "" {
		panic(err_jwt_not_set)
	}

	if os.Getenv("BACKEND_PORT") != "" {
		port = os.Getenv("BACKEND_PORT")
	}

	gin.SetMode(gin.TestMode)

	serverConfig := &srv.ServerConfig{
		Port:    port,
		Testing: true,
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
	}

	server, err = srv.New(serverConfig)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	server.InitHandler()
	code := m.Run()
	os.Exit(code)
}

func TestAuth(t *testing.T) {
	// TODO: Fix this test (and endpoint)

	// This should set the default admin username and password to the first ones that were provided
	// authBody := srv.AuthBody{
	// 	Username: "admin",
	// 	Password: "admin",
	// }

	// body, _ := json.Marshal(authBody)

	// req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(body))
	// req.Header.Set("Content-Type", "application/json")

	// w := httptest.NewRecorder()
	// server.Router.ServeHTTP(w, req)

	// var response map[string]interface{}
	// err := json.Unmarshal(w.Body.Bytes(), &response)

	// assert.NoError(t, err, "Failed to unmarshal response: %v", err)

	// _, exists := response["token"]

	// assert.True(t, exists, "Token should be present in the response body, but it was not found. %v", response)
}

func TestGalleryGroupCreation(t *testing.T) {
	name := "testing_group"

	req, _ := http.NewRequest("POST", "/gallery/groups/new/"+name, nil)

	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, err, "Failed to unmarshal response: %v", err)

	// Make another request and expect it to fail because of a duped name
	http.NewRequest("POST", "/gallery/groups/new/"+name, nil)

	w2 := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w2.Code)
}
