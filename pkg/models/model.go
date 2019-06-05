package models


import (
	"errors"
	"time"
)

// ErrNoRecord private variable for model use
var ErrNoRecord = errors.New("models: no matching record found")

//Snippet model entity struct
type Snippet struct {
	ID		int
	Title	string
	Content	string
	Created	time.Time
	Expires	time.Time
}

// Snip definition dto model 
type Snip struct {
	ID		int			`json:"id"`
	Title	string		`json:"title"`
	Content	string		`json:"content"`
	Created	string		`json:"create"`
	Expires	string		`json:"expires"`
}