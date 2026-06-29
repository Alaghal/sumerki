package domain

import "time"

type Mission struct {
	ID          string
	KingdomID   string
	Key         string
	Type        string
	Status      string
	StartedAt   time.Time
	FinishesAt  time.Time
	CompletedAt *time.Time
	ResultJSON  []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MissionUnit struct {
	ID             string
	MissionID      string
	UnitType       string
	AmountSent     int64
	AmountLost     int64
	AmountReturned int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
