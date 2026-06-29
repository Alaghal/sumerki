package domain

import "time"

type Resources struct {
	KingdomID        string
	Gold             int64
	Food             int64
	Wood             int64
	Stone            int64
	Population       int64
	LastCalculatedAt time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
