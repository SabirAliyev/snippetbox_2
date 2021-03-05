package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"sabiraliyev.net/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	// Bcrypt hash of plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Calling the Begin() method on the connection pool creates a new sql.Tx
	// object, which represents the in-progress database transaction.
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// Call Exec() on the transaction, passing in your statement and any parameters.
	// It`s important to notice that tx.Exec() is called on the transaction object,
	// just created, NOT the connection pool. Although we we`re using tx.Exec() here
	// you can also use tx.Query() and tx.QueryRow() in exactly the same way.
	insertStmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())`
	user, err := m.GetUserByEmail(email)
	if err == nil {
		fmt.Println("User ID (2): ", user.ID)
		if user.ID != 0 {
			return models.ErrDuplicateEmail
		} else {
			_, err = tx.Exec(insertStmt, name, email, string(hashedPassword))
			err = tx.Commit()
		}
	} else {
		fmt.Println("Error: ", err)
	}

	rollBackErr := tx.Rollback()
	if rollBackErr != nil {
		//TODO: Handle Rollback error.
	}
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// retrieve the id and hashed password associated with given email. If no matching email exist,
	// or the user is not active, we return theErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = $1 AND active = TRUE`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don`t, we return theErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// Verify whether a user with provided email address and password. Returns the relevant user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT  id, name, email, created, active, administrator FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active, &u.Administrator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (m *UserModel) GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT  id, name, email, created, active, administrator FROM users WHERE email = $1`
	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active, &u.Administrator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}
