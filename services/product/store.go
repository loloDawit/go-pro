package product

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
)

type ProductStore struct {
	db  *sql.DB
	cfg *config.Config
}

func NewProductStore(db *sql.DB, cfg *config.Config) *ProductStore {
	return &ProductStore{db: db, cfg: cfg}
}

func (s *ProductStore) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT id, name, description, image, price, quantity, createdAt FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []types.Product{}
	for rows.Next() {
		p := types.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *ProductStore) GetProductByID(id int) (*types.Product, error) {
	row := s.db.QueryRow("SELECT * FROM products WHERE id = $1", id)

	p := new(types.Product)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProductStore) CreateProduct(p types.Product) (int, error) {
	var newID int
	err := s.db.QueryRow(
		"INSERT INTO products (name, description, image, price, quantity) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		p.Name, p.Description, p.Image, p.Price, p.Quantity,
	).Scan(&newID)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		return 0, err
	}

	return newID, nil
}

func (s *ProductStore) UpdateProductQuantityWithTransaction(p types.Product) error {
	// Begin a new transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Ensure the transaction is either committed or rolled back
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() // Rollback in case of panic
			panic(p)
		} else if err != nil {
			tx.Rollback() // Rollback in case of error
		} else {
			err = tx.Commit() // Commit if no errors
		}
	}()

	// Retrieve the initial quantity within the transaction
	var initialQuantity int
	err = tx.QueryRow("SELECT quantity FROM products WHERE id = $1 FOR UPDATE", p.ID).Scan(&initialQuantity)
	if err != nil {
		fmt.Printf("Error retrieving initial quantity: %v\n", err)
		return err
	}
	fmt.Printf("Initial quantity for product ID %d: %d\n", p.ID, initialQuantity)

	// Execute the SQL update statement within the transaction
	result, err := tx.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", p.Quantity, p.ID)
	if err != nil {
		fmt.Printf("Error executing update: %v\n", err)
		return err
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting affected rows: %v\n", err)
		return err
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated, check if the product ID exists and quantity is valid")
	}

	// Retrieve and log the updated quantity within the transaction
	var updatedQuantity int
	err = tx.QueryRow("SELECT quantity FROM products WHERE id = $1", p.ID).Scan(&updatedQuantity)
	if err != nil {
		fmt.Printf("Error retrieving updated quantity: %v\n", err)
		return err
	}
	fmt.Printf("Updated product ID: %d, New quantity: %d\n", p.ID, updatedQuantity)

	// Verify the update
	expectedQuantity := initialQuantity - p.Quantity
	if updatedQuantity != expectedQuantity {
		return fmt.Errorf("quantity update verification failed: expected %d but got %d", expectedQuantity, updatedQuantity)
	}

	return nil
}
