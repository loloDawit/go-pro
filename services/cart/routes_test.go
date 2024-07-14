package cart

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
	"github.com/stretchr/testify/assert"
)

type mockOrderStore struct {
	CreateOrderFunc     func(order types.Order) (int, error)
	CreateOrderItemFunc func(item types.OrderItem) error
}

func (m *mockOrderStore) CreateOrder(order types.Order) (int, error) {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(order)
	}
	return 0, nil
}

func (m *mockOrderStore) CreateOrderItem(item types.OrderItem) error {
	if m.CreateOrderItemFunc != nil {
		return m.CreateOrderItemFunc(item)
	}
	return nil
}

type mockProductStore struct {
	GetProductsFunc                          func() ([]types.Product, error)
	GetProductByIDFunc                       func(id int) (*types.Product, error)
	UpdateProductQuantityWithTransactionFunc func(product types.Product) error
	CreateProductFunc                        func(product types.Product) (int, error)
}

func (m *mockProductStore) GetProducts() ([]types.Product, error) {
	if m.GetProductsFunc != nil {
		return m.GetProductsFunc()
	}
	return nil, nil
}

func (m *mockProductStore) GetProductByID(id int) (*types.Product, error) {
	if m.GetProductByIDFunc != nil {
		return m.GetProductByIDFunc(id)
	}
	return nil, nil
}

func (m *mockProductStore) UpdateProductQuantityWithTransaction(product types.Product) error {
	if m.UpdateProductQuantityWithTransactionFunc != nil {
		return m.UpdateProductQuantityWithTransactionFunc(product)
	}
	return nil
}

func (m *mockProductStore) CreateProduct(product types.Product) (int, error) {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(product)
	}
	return 0, nil
}

func generateTestToken(secret []byte, userID int, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(userID),
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TestCheckoutRoute(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "testsecret",
		},
	}

	mockOrderStore := &mockOrderStore{
		CreateOrderFunc: func(order types.Order) (int, error) {
			return 123, nil
		},
		CreateOrderItemFunc: func(item types.OrderItem) error {
			return nil
		},
	}

	mockProductStore := &mockProductStore{
		GetProductByIDFunc: func(id int) (*types.Product, error) {
			return &types.Product{
				ID:       id,
				Name:     "Test Product",
				Price:    10,
				Quantity: 100,
			}, nil
		},
		UpdateProductQuantityWithTransactionFunc: func(product types.Product) error {
			return nil
		},
	}

	handler := NewHandlers(mockOrderStore, mockProductStore, cfg)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	payload := types.CartCheckoutPayload{
		Items: []types.CartItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/cart/checkout", bytes.NewReader(body))
	assert.NoError(t, err)

	// Generate a test token and add it to the request header
	token, err := generateTestToken([]byte(cfg.JWT.Secret), 1, time.Hour)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Debugging: Print the response for troubleshooting
	if rr.Code != http.StatusOK {
		t.Logf("Response Code: %v", rr.Code)
		t.Logf("Response Body: %v", rr.Body.String())
	}

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualResponse types.CreateOrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expectedResponse := types.CreateOrderResponse{
		ID:      123,
		Total:   20.0,
		Message: "Order created successfully",
	}
	assert.Equal(t, expectedResponse, actualResponse)
}
