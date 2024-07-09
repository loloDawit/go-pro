package product

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewProductStore(db, cfg)

	tests := []struct {
		name         string
		mockQuery    func()
		expectedProducts []types.Product
		expectedErr  error
	}{
		{
			name: "Products found",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "image", "price", "quantity", "createdAt"}).
					AddRow(1, "Product 1", "Description 1", "image1.jpg", 10.5, 100, time.Now()).
					AddRow(2, "Product 2", "Description 2", "image2.jpg", 20.0, 200, time.Now())
				mock.ExpectQuery("SELECT id, name, description, image, price, quantity, createdAt FROM products").
					WillReturnRows(rows)
			},
			expectedProducts: []types.Product{
				{ID: 1, Name: "Product 1", Description: "Description 1", Image: "image1.jpg", Price: 10.5, Quantity: 100},
				{ID: 2, Name: "Product 2", Description: "Description 2", Image: "image2.jpg", Price: 20.0, Quantity: 200},
			},
			expectedErr: nil,
		},
		{
			name: "Database error",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, name, description, image, price, quantity, createdAt FROM products").
					WillReturnError(sql.ErrConnDone)
			},
			expectedProducts: nil,
			expectedErr:      sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			products, err := store.GetProducts()
			assert.Equal(t, tt.expectedErr, err)
			for i, product := range products {
				assert.Equal(t, tt.expectedProducts[i].ID, product.ID)
				assert.Equal(t, tt.expectedProducts[i].Name, product.Name)
				assert.Equal(t, tt.expectedProducts[i].Description, product.Description)
				assert.Equal(t, tt.expectedProducts[i].Image, product.Image)
				assert.Equal(t, tt.expectedProducts[i].Price, product.Price)
				assert.Equal(t, tt.expectedProducts[i].Quantity, product.Quantity)
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewProductStore(db, cfg)

	tests := []struct {
		name         string
		id           int
		mockQuery    func()
		expectedProduct *types.Product
		expectedErr  error
	}{
		{
			name: "Product found",
			id:   1,
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "image", "price", "quantity", "createdAt"}).
					AddRow(1, "Product 1", "Description 1", "image1.jpg", 10.5, 100, time.Now())
				mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedProduct: &types.Product{ID: 1, Name: "Product 1", Description: "Description 1", Image: "image1.jpg", Price: 10.5, Quantity: 100},
			expectedErr: nil,
		},
		{
			name: "Product not found",
			id:   1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedProduct: nil,
			expectedErr:     sql.ErrNoRows,
		},
		{
			name: "Database error",
			id:   1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedProduct: nil,
			expectedErr:     sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			product, err := store.GetProductByID(tt.id)
			assert.Equal(t, tt.expectedErr, err)
			if product != nil {
				assert.Equal(t, tt.expectedProduct.ID, product.ID)
				assert.Equal(t, tt.expectedProduct.Name, product.Name)
				assert.Equal(t, tt.expectedProduct.Description, product.Description)
				assert.Equal(t, tt.expectedProduct.Image, product.Image)
				assert.Equal(t, tt.expectedProduct.Price, product.Price)
				assert.Equal(t, tt.expectedProduct.Quantity, product.Quantity)
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewProductStore(db, cfg)

	tests := []struct {
		name        string
		product     types.Product
		mockExec    func()
		expectedID  int
		expectedErr error
	}{
		{
			name: "Successful creation",
			product: types.Product{
				Name:        "Product 1",
				Description: "Description 1",
				Image:       "image1.jpg",
				Price:       10.5,
				Quantity:    100,
			},
			mockExec: func() {
				mock.ExpectQuery("INSERT INTO products \\(name, description, image, price, quantity\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
					WithArgs("Product 1", "Description 1", "image1.jpg", 10.5, 100).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Database error",
			product: types.Product{
				Name:        "Product 1",
				Description: "Description 1",
				Image:       "image1.jpg",
				Price:       10.5,
				Quantity:    100,
			},
			mockExec: func() {
				mock.ExpectQuery("INSERT INTO products \\(name, description, image, price, quantity\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
					WithArgs("Product 1", "Description 1", "image1.jpg", 10.5, 100).
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:  0,
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExec()
			id, err := store.CreateProduct(tt.product)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedID, id)
		})
	}
}

func TestUpdateProductQuantityWithTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewProductStore(db, cfg)

	tests := []struct {
		name        string
		product     types.Product
		mockQuery   func()
		expectedErr error
	}{
		{
			name: "Successful update",
			product: types.Product{
				ID:       1,
				Quantity: 10,
			},
			mockQuery: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT quantity FROM products WHERE id = \\$1 FOR UPDATE").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(20))
				mock.ExpectExec("UPDATE products SET quantity = quantity - \\$1 WHERE id = \\$2").
					WithArgs(10, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery("SELECT quantity FROM products WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(10))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "Product not found",
			product: types.Product{
				ID:       1,
				Quantity: 10,
			},
			mockQuery: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT quantity FROM products WHERE id = \\$1 FOR UPDATE").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "Database error during update",
			product: types.Product{
				ID:       1,
				Quantity: 10,
			},
			mockQuery: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT quantity FROM products WHERE id = \\$1 FOR UPDATE").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(20))
				mock.ExpectExec("UPDATE products SET quantity = quantity - \\$1 WHERE id = \\$2").
					WithArgs(10, 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			err := store.UpdateProductQuantityWithTransaction(tt.product)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
