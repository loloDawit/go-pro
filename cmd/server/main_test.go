package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/loloDawit/ecom/config"
)

func setupTestEnv(t *testing.T) (*APIServer, sqlmock.Sqlmock) {
	t.Helper()

	// Mock configuration
	cfg := &config.Config{
		DBuser:     "testuser",
		DBpassword: "testpassword",
		DBaddr:     "localhost",
		DBname:     "testdb",
		Address:    ":8080",
	}

	// Initialize the mock database with MonitorPingsOption
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to open mock sql db, %v", err)
	}

	server := NewAPIServer(cfg.Address, db, cfg)

	return server, mock
}

func TestHealthCheckHandler(t *testing.T) {
	server, mock := setupTestEnv(t)
	defer server.db.Close()

	// Expect a ping to the database
	mock.ExpectPing().WillReturnError(nil)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	server.healthCheckHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", res.Status)
	}
}

func TestServerStart(t *testing.T) {
	server, mock := setupTestEnv(t)
	defer server.db.Close()

	// Expect a ping to the database
	mock.ExpectPing().WillReturnError(nil)

	// Run the server in a separate goroutine
	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	// Give the server some time to start
	<-time.After(time.Second)

	// Test the /health endpoint
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to send request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
}


func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Exit with the result of the tests
	os.Exit(code)
}
