package models


import (
	"errors"
	"time"
)

// ErrNoRecord private variable for model use
var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

//Snippet model entity struct
type Snippet struct {
	ID		int
	Title	string
	Content	string
	Created	time.Time
	Expires	time.Time
}

type User struct {
	ID				int
	Name			string
	Email			string
	HashedPassword	[]byte
	Created			time.Time
	Active			bool
}

// Snip definition dto model json option
type Snip struct {
	ID		int			`json:"id"`
	Title	string		`json:"title"`
	Content	string		`json:"content"`
	Created	string		`json:"create"`
	Expires	string		`json:"expires"`
}