package product

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/types"
	"github.com/stretchr/testify/assert"
)

type mockProductStore struct {
	GetProductsFunc                          func() ([]types.Product, error)
	GetProductByIDFunc                       func(id int) (*types.Product, error)
	CreateProductFunc                        func(product types.Product) (int, error)
	UpdateProductQuantityWithTransactionFunc func(product types.Product) error
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

func (m *mockProductStore) CreateProduct(product types.Product) (int, error) {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(product)
	}
	return 0, nil
}

func (m *mockProductStore) UpdateProductQuantityWithTransaction(product types.Product) error {
	if m.UpdateProductQuantityWithTransactionFunc != nil {
		return m.UpdateProductQuantityWithTransactionFunc(product)
	}
	return nil
}

func TestGetProductsRoute(t *testing.T) {
	mockStore := &mockProductStore{
		GetProductsFunc: func() ([]types.Product, error) {
			return []types.Product{
				{Name: "Product1"},
				{Name: "Product2"},
			}, nil
		},
	}

	handler := NewHandlers(mockStore)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	req, err := http.NewRequest("GET", "/products", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response body to check specific fields
	var actualProducts []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualProducts)
	assert.NoError(t, err)

	expectedNames := []string{"Product1", "Product2"}
	for i, product := range actualProducts {
		assert.Equal(t, expectedNames[i], product["name"])
	}
}

func TestCreateProductRoute(t *testing.T) {
	mockStore := &mockProductStore{
		CreateProductFunc: func(product types.Product) (int, error) {
			return 123, nil
		},
	}

	handler := NewHandlers(mockStore)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	payload := types.CreateProductPayload{
		Name:        "Test Product",
		Description: "Test Description",
		Image:       "test.jpg",
		Price:       100,
		Quantity:    10,
	}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products", bytes.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := `{"id":123,"message":"Product created successfully"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetProductRoute(t *testing.T) {
	mockStore := &mockProductStore{
		GetProductByIDFunc: func(id int) (*types.Product, error) {
			return &types.Product{Name: "Product1"}, nil
		},
	}

	handler := NewHandlers(mockStore)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	req, err := http.NewRequest("GET", "/products/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response body to check specific fields
	var actualProduct map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualProduct)
	assert.NoError(t, err)

	assert.Equal(t, "Product1", actualProduct["name"])
}
