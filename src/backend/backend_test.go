package main

import (
	"bash06/strona-fundacja/src/backend/core"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	baseUrl = "http://localhost:2333"
)

var (
	jwtToken       string
	postStr        = []byte(`{"title": "Test title", "content": "Test content"}`)
	credentialsStr = []byte(`{"username": "admin", "password": "admin"}`)
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode) // Disable debug gin logs
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	db := core.Database{
		Memory: true,
	}
	db.Init()

	r := setupRouter(&db, true)
	startServer(r)

	time.Sleep(1 * time.Second)

	code := m.Run()

	stopServer()
	os.Exit(code)
}

// This will most likely never fail but better safe than sorry
func TestAdminLogin(t *testing.T) {
	mockCredentials := bytes.NewBuffer(credentialsStr)

	resp, err := http.Post(baseUrl+"/admin/login", "application/json", mockCredentials)
	if err != nil {
		t.Fatalf("Failed to login as test admin user: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got %v", resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Response body: %s", bodyBytes)
	}

	defer resp.Body.Close()

	var body struct {
		Token string `json:"token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	jwtToken = body.Token

}

func TestCreatePost(t *testing.T) {
	mockPost := bytes.NewBuffer(postStr)

	client := &http.Client{}

	req, err := http.NewRequest("POST", baseUrl+"/blog/post", mockPost)
	if err != nil {
		t.Errorf("Failed to prepare create post request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to create post record: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %v", resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Response body: %s", bodyBytes)
	}
}

func TestGetPost(t *testing.T) {
	resp, err := http.Get(baseUrl + "/blog/post/1")
	if err != nil {
		t.Errorf("Failed to get post: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %v", resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Response body: %s", bodyBytes)
	}
}

func TestEditPost(t *testing.T) {
	req, err := http.NewRequest("DELETE", baseUrl+"/blog/post/1", nil)
	if err != nil {
		t.Errorf("Prepare request failed: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", jwtToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %v", resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Response body: %s", bodyBytes)
	}
}
