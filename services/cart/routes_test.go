package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/services/auth"
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

func MockJWTMiddleware(secret []byte) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := strconv.Itoa(1)
			ctx := context.WithValue(r.Context(), types.UserIDKey, userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func BypassJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func MockEmptyUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate an authorized request by setting an empty context
		ctx := context.WithValue(r.Context(), types.UserIDKey, "")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}



func TestCheckoutRoute(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "testsecret",
		},
	}

	tests := []struct {
		name                 string
		payload              interface{}
		mockOrderStore       *mockOrderStore
		mockProductStore     *mockProductStore
		expectedStatus       int
		expectedResponseBody string
		bypassMiddleware     bool
	}{
		{
			name:    "Invalid Payload",
			payload: "invalid json",
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid payload"}`,
		},
		{
			name:    "Validation Errors",
			payload: types.CartCheckoutPayload{},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid payload: Key: 'CartCheckoutPayload.Items' Error:Field validation for 'Items' failed on the 'required' tag"}`,
		},
		{
			name: "Empty Cart",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Cart is empty"}`,
		},
		{
			name: "Product Not Found",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
				GetProductByIDFunc: func(id int) (*types.Product, error) {
					return nil, http.ErrNotSupported
				},
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Product not found"}`,
		},
		{
			name: "Out of Stock",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
				GetProductByIDFunc: func(id int) (*types.Product, error) {
					return &types.Product{
						ID:       id,
						Name:     "Test Product",
						Price:    10,
						Quantity: 0,
					}, nil
				},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Product Test Product is out of stock"}`,
		},
		{
			name: "Insufficient Quantity",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 10},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
				GetProductByIDFunc: func(id int) (*types.Product, error) {
					return &types.Product{
						ID:       id,
						Name:     "Test Product",
						Price:    10,
						Quantity: 5,
					}, nil
				},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Product Test Product has only 5 items left"}`,
		},
		{
			name: "Product Update Failure",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
				GetProductByIDFunc: func(id int) (*types.Product, error) {
					return &types.Product{
						ID:       id,
						Name:     "Test Product",
						Price:    10,
						Quantity: 100,
					}, nil
				},
				UpdateProductQuantityWithTransactionFunc: func(product types.Product) error {
					return http.ErrNotSupported
				},
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to update product quantity"}`,
		},
		{
			name: "Order Creation Failure",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 0, http.ErrNotSupported
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
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
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to create order"}`,
		},
		{
			name: "Order Item Creation Failure",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return http.ErrNotSupported
				},
			},
			mockProductStore: &mockProductStore{
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
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Failed to create order item"}`,
		},
		{
			name: "Successful Checkout",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
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
			},
			expectedStatus:       http.StatusOK,
			expectedResponseBody: `{"id":123,"total":20,"message":"Order created successfully"}`,
		},
		{
			name: "Error Getting User ID from Context",
			payload: types.CartCheckoutPayload{
				Items: []types.CartItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			mockOrderStore: &mockOrderStore{
				CreateOrderFunc: func(order types.Order) (int, error) {
					return 123, nil
				},
				CreateOrderItemFunc: func(item types.OrderItem) error {
					return nil
				},
			},
			mockProductStore: &mockProductStore{
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
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedResponseBody: `{"error":"user ID not found in context"}`,
			bypassMiddleware:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandlers(tt.mockOrderStore, tt.mockProductStore, cfg)

			// Create a router instance without middleware for this specific test case
			var router *mux.Router
			if tt.bypassMiddleware {
				router = mux.NewRouter()
			} else {
				router = mux.NewRouter()
				router.Use(MockJWTMiddleware([]byte(cfg.JWT.Secret)))
			}

			handler.RegisterRoutes(router)

			body, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/cart/checkout", bytes.NewReader(body))
			assert.NoError(t, err)

			// Generate a test token and add it to the request header
			token, err := generateTestToken([]byte(cfg.JWT.Secret), 1, time.Hour)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			if tt.bypassMiddleware {
				// Directly set the context without the user ID
				req = req.WithContext(context.Background())
				log.Println("Removed User ID from context for test:", tt.name)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			log.Println("Test:", tt.name)
			log.Println("Expected Status:", tt.expectedStatus, "Actual Status:", rr.Code)
			log.Println("Expected Body:", tt.expectedResponseBody, "Actual Body:", rr.Body.String())

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestCheckoutRouteErrorGettingUserID(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "testsecret",
		},
	}

	handler := NewHandlers(&mockOrderStore{
		CreateOrderFunc: func(order types.Order) (int, error) {
			return 123, nil
		},
		CreateOrderItemFunc: func(item types.OrderItem) error {
			return nil
		},
	}, &mockProductStore{
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
	}, cfg)

	// Create a router instance with the middleware
	router := mux.NewRouter()
	router.HandleFunc("/cart/checkout", auth.JWTMiddleware([]byte(cfg.JWT.Secret))(handler.checkout)).Methods("POST")

	body, err := json.Marshal(types.CartCheckoutPayload{
		Items: []types.CartItem{
			{ProductID: 1, Quantity: 2},
		},
	})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/cart/checkout", bytes.NewReader(body))
	assert.NoError(t, err)

	// Generate a test token and add it to the request header
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": "1",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Set the special header to bypass setting the user ID in the middleware
	req.Header.Set("X-Bypass-UserID", "true")
	log.Println("Set special header to bypass setting the user ID for test: Error Getting User ID from Context")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	log.Println("Test: Error Getting User ID from Context")
	log.Println("Expected Status: 500, Actual Status:", rr.Code)
	log.Println("Expected Body: {\"error\":\"user ID not found in context\"}, Actual Body:", rr.Body.String())

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{"error":"user ID not found in context"}`, rr.Body.String())
}