package domain

import "time"

type PatronRelation struct {
	ID        string
	KingdomID string
	Patron    string
	Favor     int
	Standing  string
	JoinedAt  time.Time
	LeftAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
