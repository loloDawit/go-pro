package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/services/auth"
	"github.com/loloDawit/ecom/types"
	"github.com/loloDawit/ecom/utils"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	store            types.UserStore
	cfg              *config.Config
	comparePasswords func(string, string) error
	generateToken    func([]byte, int, time.Duration) (string, error)
}

func NewHandlers(store types.UserStore, cfg *config.Config) *Handler {
	return &Handler{
		store:            store,
		cfg:              cfg,
		comparePasswords: auth.ComparePasswords,
		generateToken:    auth.GenerateToken,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/signup", h.signUp).Methods("POST")
	r.HandleFunc("/login", h.login).Methods("POST")

}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidRequestBody)
		return
	}
	var payload types.SignupUserPayload

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

	// check if the user already exists
	if err := h.checkUserExists(payload.Email); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		log.Printf("%s: %v", utils.ErrHashingPassword, err)
		utils.WriteError(w, http.StatusInternalServerError, utils.ErrInternalServerError)
		return
	}

	// if the user does not exist, create the user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Printf("%s: %v", utils.ErrCreatingUser, err)
		utils.WriteError(w, http.StatusInternalServerError, utils.ErrInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidRequestBody)
		return
	}

	var payload types.LoginUserPayload
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

	// get the user by email
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, utils.ErrUserNotFound)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, utils.ErrInternalServerError)
		return
	}

	// compare the password
	if err := h.comparePasswords(user.Password, payload.Password); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, utils.ErrUnauthorized)
		return
	}

	// generate a token
	expiration := time.Second * time.Duration(h.cfg.JWT.Expiration)
	token, err := h.generateToken([]byte(h.cfg.JWT.Secret), user.ID, expiration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, utils.ErrInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// checkUserExists checks if a user with the given email already exists
func (h *Handler) checkUserExists(email string) error {
	_, err := h.store.GetUserByEmail(email)
	if err == nil {
		return fmt.Errorf(utils.ErrUserAlreadyExists)
	}
	if err == sql.ErrNoRows {
		// User not found, proceed
		return nil
	}

	log.Printf("error checking user existence: %v", err)
	return fmt.Errorf(utils.ErrInternalServerError)
}
