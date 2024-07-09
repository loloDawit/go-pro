package user

import (
	"database/sql"
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

	t.Run("User found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password"}).
			AddRow(1, "John", "Doe", "john.doe@example.com", "hashedpassword")
		mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE email = \\$1").
			WithArgs("john.doe@example.com").
			WillReturnRows(rows)

		user, err := store.GetUserByEmail("john.doe@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE email = \\$1").
			WithArgs("john.doe@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := store.GetUserByEmail("john.doe@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewUserStore(db, cfg)

	mock.ExpectExec("INSERT INTO users \\(firstName, lastName, email, password\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs("John", "Doe", "john.doe@example.com", "hashedpassword").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.CreateUser(types.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashedpassword",
	})
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &config.Config{}
	store := NewUserStore(db, cfg)

	t.Run("User found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password"}).
			AddRow(1, "John", "Doe", "john.doe@example.com", "hashedpassword")
		mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(rows)

		user, err := store.GetUserByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, firstName, lastName, email, password FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		user, err := store.GetUserByID(1)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}
