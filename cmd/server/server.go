package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/config"
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
