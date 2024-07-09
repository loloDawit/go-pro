package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/db"
	"github.com/loloDawit/ecom/services/cart"
	"github.com/loloDawit/ecom/services/order"
	"github.com/loloDawit/ecom/services/product"
	"github.com/loloDawit/ecom/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
	cfg  *config.Config
}

func NewAPIServer(addr string, db *sql.DB, cfg *config.Config) *APIServer {
	return &APIServer{addr: addr, db: db, cfg: cfg}
}

func (s *APIServer) Start() error {
	// initialize the router
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// initialize the user handler
	userHandler := user.NewHandlers(user.NewUserStore(s.db, s.cfg), s.cfg)
	userHandler.RegisterRoutes(subrouter)

	// initialize the product handler
	productHandler := product.NewHandlers(product.NewProductStore(s.db, s.cfg))
	productHandler.RegisterRoutes(subrouter)

	// initialize the cart handler
	orderStore := order.NewOrderStore(s.db, s.cfg)
	cartHandler := cart.NewHandlers(orderStore, product.NewProductStore(s.db, s.cfg), s.cfg)
	cartHandler.RegisterRoutes(subrouter)

	// add log for server listening
	log.Printf("Server is listening on %s", s.addr)

	return http.ListenAndServe(s.addr, router)
}

func main() {
	// Check current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	log.Println("Current working directory:", cwd)

	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatalf(".env file does not exist")
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Ensure CONFIG_DIRECTORY is set
	if os.Getenv("CONFIG_DIRECTORY") == "" {
		os.Setenv("CONFIG_DIRECTORY", "./config")
	}

	// Initialize configuration
	ctx := context.Background()
	directory := os.Getenv("CONFIG_DIRECTORY")
	environment := os.Getenv("ENV")
	deployment := os.Getenv("DEPLOYMENT")

	cfg := config.LoadConfig(ctx, directory, environment, deployment)

	fmt.Println(cfg)

	// Initialize the database
	// Construct the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", cfg.DBuser, cfg.DBpassword, cfg.DBaddr, cfg.DBname)

	fmt.Println(connStr)

	db, err := db.NewSQLDatabase(connStr)
	if err != nil {
		panic(err)
	}

	initStorage(db)

	// Initialize and start the server
	server := NewAPIServer(cfg.Address, db, cfg)
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	log.Println("Database is connected")
}
