package user

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewUserStore(db, cfg)

	tests := []struct {
		name         string
		email        string
		mockQuery    func()
		expectedUser *types.User
		expectedErr  error
	}{
		{
			name:  "User found",
			email: "john.doe@example.com",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password"}).
					AddRow(1, "John", "Doe", "john.doe@example.com", "hashedpassword")
				mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE email = \\$1").
					WithArgs("john.doe@example.com").
					WillReturnRows(rows)
			},
			expectedUser: &types.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashedpassword",
			},
			expectedErr: nil,
		},
		{
			name:  "User not found",
			email: "john.doe@example.com",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE email = \\$1").
					WithArgs("john.doe@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser: nil,
			expectedErr:  sql.ErrNoRows,
		},
		{
			name:  "Database error",
			email: "john.doe@example.com",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE email = \\$1").
					WithArgs("john.doe@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser: nil,
			expectedErr:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			user, err := store.GetUserByEmail(tt.email)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, user)
		})
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewUserStore(db, cfg)

	tests := []struct {
		name        string
		user        types.User
		mockExec    func()
		expectedErr error
	}{
		{
			name: "Successful creation",
			user: types.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashedpassword",
			},
			mockExec: func() {
				mock.ExpectExec("INSERT INTO users \\(firstName, lastName, email, password\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
					WithArgs("John", "Doe", "john.doe@example.com", "hashedpassword").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Database error",
			user: types.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashedpassword",
			},
			mockExec: func() {
				mock.ExpectExec("INSERT INTO users \\(firstName, lastName, email, password\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
					WithArgs("John", "Doe", "john.doe@example.com", "hashedpassword").
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExec()
			err := store.CreateUser(tt.user)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewUserStore(db, cfg)

	tests := []struct {
		name         string
		id           int
		mockQuery    func()
		expectedUser *types.User
		expectedErr  error
	}{
		{
			name: "User found",
			id:   1,
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password"}).
					AddRow(1, "John", "Doe", "john.doe@example.com", "hashedpassword")
				mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser: &types.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashedpassword",
			},
			expectedErr: nil,
		},
		{
			name: "User not found",
			id:   1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser: nil,
			expectedErr:  fmt.Errorf("user not found"),
		},
		{
			name: "Database error",
			id:   1,
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser: nil,
			expectedErr:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			user, err := store.GetUserByID(tt.id)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, user)
		})
	}
}
