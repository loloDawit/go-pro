package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
)

type mockUserStore struct {
	db *sql.DB
}

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	config := &config.Config{}
	// create a new user service
	handler := NewHandlers(userStore, config)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		// create a new user
		u := types.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "da",
			Password:  "password",
		}
		marshalled, _ := json.Marshal(u)

		req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/signup", handler.signUp)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	// happy path
	t.Run("should create a new user", func(t *testing.T) {
		u := types.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "da@test.com",
			Password:  "password",
		}
		marshalled, _ := json.Marshal(u)

		req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/signup", handler.signUp)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

	})

}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

// CreateUser implements types.UserStore.
func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}

// GetUserByID implements types.UserStore.
func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}
