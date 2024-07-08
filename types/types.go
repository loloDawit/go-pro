package types

import "time"

type contextKey string

const UserIDKey contextKey = "userID"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
	GetUserByID(id int) (*User, error)
}

type SignupUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=20"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	Image       string    `json:"image"`
	Quantity    int       `json:"quantity"`
}

type ProductStore interface {
	GetProductByID(id int) (*Product, error)
	GetProducts() ([]Product, error)
	CreateProduct(Product) (int, error)
	UpdateProductQuantityWithTransaction(Product) error
}

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type CreateProductResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// cart
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderId"`
	ProductID int       `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}

type CreateOrderResponse struct {
	ID      int     `json:"id"`
	Total   float64 `json:"total"`
	Message string  `json:"message"`
}
