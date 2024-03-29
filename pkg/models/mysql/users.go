package mysql

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"

	// "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	// "strings"

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

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())`

	// The Exec() method to insert the user details and hashed password into the users table.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check whether the error
		// has the type *mysql.MySQLError. If it does, the error will be assigned to the MySQLError
		// variable. We can then check whether or not the error relates to our users_uc_email key by
		// checking the contents of the message string. If it does, we return an ErrDuplicateEmail error.
		pqError := err.(*pq.Error)
		if errors.As(err, &pqError) {
			//if pqError.Code.Name() == 1062 && strings.Contains(pqError.Message, "users_uc_email") {
			//	return models.ErrDuplicvateEmail
			//}
		}
		return err
	}

	return nil
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
