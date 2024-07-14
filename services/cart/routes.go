package cart

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/services/auth"
	"github.com/loloDawit/ecom/types"
	"github.com/loloDawit/ecom/utils"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	cfg          *config.Config
}

func NewHandlers(store types.OrderStore, productStore types.ProductStore, cfg *config.Config) *Handler {
	return &Handler{store: store, productStore: productStore, cfg: cfg}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cart/checkout", auth.JWTMiddleware([]byte(h.cfg.JWT.Secret))(h.checkout)).Methods("POST")
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		fmt.Println("Error getting user ID from context:", err)
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("User ID:", userID)

	var cartPayload types.CartCheckoutPayload
	err = utils.ReadJSON(r, &cartPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := utils.Validate.Struct(cartPayload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", validationErrors))
		return
	}

	if len(cartPayload.Items) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Cart is empty")
		return
	}

	var totalPrice float64
	for _, item := range cartPayload.Items {
		product, err := h.productStore.GetProductByID(item.ProductID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Product not found")
			return
		}

		if product.Quantity <= 0 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Product %s is out of stock", product.Name))
			return
		}

		if item.Quantity > product.Quantity {
			utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Product %s has only %d items left", product.Name, product.Quantity))
			return
		}

		totalPrice += float64(product.Price) * float64(item.Quantity)

		err = h.productStore.UpdateProductQuantityWithTransaction(types.Product{
			ID:       product.ID,
			Quantity: item.Quantity,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update product quantity")
			return
		}
	}

	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Seattle, WA",
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// create the order item
	for _, item := range cartPayload.Items {
		err = h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     totalPrice,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to create order item")
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, types.CreateOrderResponse{
		ID:      orderID,
		Total:   totalPrice,
		Message: "Order created successfully",
	})
}

func getUserIDFromContext(ctx context.Context) (int, error) {
	userIDStr, ok := ctx.Value(types.UserIDKey).(string)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return strconv.Atoi(userIDStr)
}
