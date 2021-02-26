package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching records found")
	// The error about incorrect email address.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// The error about duplicate emails.
	ErrDuplicvateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Message struct {
	MessageID int
	UserId    int
	Content   string
	Created   time.Time
	Expires   time.Time
	Deleted   bool
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
	Administrator  bool
}
