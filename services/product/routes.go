package product

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/types"
	"github.com/loloDawit/ecom/utils"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	store types.ProductStore
}

func NewHandlers(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/products", h.getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", h.getProduct).Methods("GET")
	r.HandleFunc("/products", h.createProduct).Methods("POST")
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.store.GetProductByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	// read the payload
	if r.Body == nil {
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidRequestBody)
		return
	}

	var payload types.CreateProductPayload
	err := utils.ReadJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidPayload)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("%s: %v", utils.ErrInvalidPayload, validationErrors))
		return
	}

	// create the product
	productID, err := h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	if err != nil {
		log.Printf("%s: %v", utils.ErrCreatingProduct, err)
		utils.WriteError(w, http.StatusInternalServerError, utils.ErrInternalServerError)
		return
	}

	response := types.CreateProductResponse{
		ID:      productID,
		Message: "Product created successfully",
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}
