package mysql

import (
	"database/sql"

	"sabiraliyev.net/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Verify whether a user with provided email address and password. Returns the relevant user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
