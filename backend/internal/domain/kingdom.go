package domain

import "time"

type Kingdom struct {
	ID        string
	UserID    string
	Name      string
	Culture   string
	Patron    *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
