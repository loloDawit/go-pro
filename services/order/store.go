package order

import (
	"database/sql"

	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
)

type OrderStore struct {
	db  *sql.DB
	cfg *config.Config
}

func NewOrderStore(db *sql.DB, cfg *config.Config) *OrderStore {
	return &OrderStore{db: db, cfg: cfg}
}

func (s *OrderStore) CreateOrder(order types.Order) (int, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO orders (userID, total, status, address) VALUES ($1, $2, $3, $4) RETURNING id",
		order.UserID, order.Total, order.Status, order.Address,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *OrderStore) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderID, productID, quantity, price) VALUES ($1, $2, $3, $4)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return err
	}

	return nil
}
