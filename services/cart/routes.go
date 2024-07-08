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
	// get the user ID from the context
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
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidPayload)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(cartPayload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("%s: %v", utils.ErrInvalidPayload, validationErrors))
		return
	}

	if len(cartPayload.Items) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Cart is empty")
		return
	}

	// get the product by id
	// loop through the cart items and get the product by id
	var totalPrice float64
	for _, item := range cartPayload.Items {
		product, err := h.productStore.GetProductByID(item.ProductID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check if the quantity is available
		if product.Quantity <= 0 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Product %s is out of stock", product.Name))
			return
		}

		// check if the quantity is not more than the available quantity
		if item.Quantity > product.Quantity {
			utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Product %s has only %d items left", product.Name, product.Quantity))
			return
		}

		// calculate the total price
		totalPrice += float64(product.Price) * float64(item.Quantity)

		// update the product quantity
		err = h.productStore.UpdateProductQuantityWithTransaction(*product)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
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
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
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
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, types.CreateOrderResponse{
		ID:      orderID,
		Total:   totalPrice,
		Message: "Order created successfully",
	})

}

// Helper function to retrieve user ID from context
func getUserIDFromContext(ctx context.Context) (int, error) {
	userIDStr, ok := ctx.Value(types.UserIDKey).(string)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return strconv.Atoi(userIDStr)
}
