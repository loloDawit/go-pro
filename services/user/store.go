package user

import (
	"database/sql"
	"fmt"

	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/types"
)

type UserStore struct {
	db  *sql.DB
	cfg *config.Config
}

func NewUserStore(db *sql.DB, cfg *config.Config) *UserStore {
	return &UserStore{db: db, cfg: cfg}
}

func (s *UserStore) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow("SELECT id, firstName, lastName, email, password FROM users WHERE email = $1", email)

	u := new(types.User)
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return u, nil
}

func (s *UserStore) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES ($1, $2, $3, $4)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow("SELECT id, firstName, lastName, email, password FROM users WHERE id = $1", id)

	u := new(types.User)
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return u, nil
}
