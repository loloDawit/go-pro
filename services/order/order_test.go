package order

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewOrderStore(db, cfg)

	tests := []struct {
		name        string
		order       types.Order
		mockQuery   func()
		expectedID  int
		expectedErr error
	}{
		{
			name: "Successful creation",
			order: types.Order{
				UserID:  1,
				Total:   100.50,
				Status:  "pending",
				Address: "123 Main St",
			},
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO orders \\(userID, total, status, address\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
					WithArgs(1, 100.50, "pending", "123 Main St").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "Database error",
			order: types.Order{
				UserID:  1,
				Total:   100.50,
				Status:  "pending",
				Address: "123 Main St",
			},
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO orders \\(userID, total, status, address\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
					WithArgs(1, 100.50, "pending", "123 Main St").
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:  0,
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			id, err := store.CreateOrder(tt.order)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedID, id)
		})
	}
}

func TestCreateOrderItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewOrderStore(db, cfg)

	tests := []struct {
		name        string
		orderItem   types.OrderItem
		mockExec    func()
		expectedErr error
	}{
		{
			name: "Successful creation",
			orderItem: types.OrderItem{
				OrderID:   1,
				ProductID: 1,
				Quantity:  2,
				Price:     50.25,
			},
			mockExec: func() {
				mock.ExpectExec("INSERT INTO order_items \\(orderID, productID, quantity, price\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
					WithArgs(1, 1, 2, 50.25).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Database error",
			orderItem: types.OrderItem{
				OrderID:   1,
				ProductID: 1,
				Quantity:  2,
				Price:     50.25,
			},
			mockExec: func() {
				mock.ExpectExec("INSERT INTO order_items \\(orderID, productID, quantity, price\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
					WithArgs(1, 1, 2, 50.25).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExec()
			err := store.CreateOrderItem(tt.orderItem)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
