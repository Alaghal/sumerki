package domain

import "time"

type Kingdom struct {
	ID        string
	UserID    string
	Name      string
	Culture   string
	Patron    *string
	Dread     int
	Honor     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
