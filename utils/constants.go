package utils

const (
	ErrInvalidRequestBody  = "please send a valid request body"
	ErrInvalidPayload      = "invalid payload"
	ErrUserAlreadyExists   = "user with this email already exists"
	ErrUserNotFound        = "user not found"
	ErrHashingPassword     = "error hashing password"
	ErrCreatingUser        = "error creating user"
	ErrInternalServerError = "internal server error"
	ErrUnauthorized        = "unauthorized"
	ErrCreatingProduct     = "error creating product"

	// success messages
	UserCreatedSuccessfully = "user created successfully"
)
