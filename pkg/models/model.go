package models


import (
	"errors"
	"time"
)

// ErrNoRecord private variable for model use
var ErrNoRecord = errors.New("models: no matching record found")

//Snippet model struct
type Snippet struct {
	ID		int
	Title	string
	Content	string
	Created	time.Time
	Expires	time.Time
}