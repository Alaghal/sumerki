package domain

import "time"

type Unit struct {
	ID        string
	KingdomID string
	Type      string
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UnitTrainingOrder struct {
	ID          string
	KingdomID   string
	UnitType    string
	Amount      int64
	Status      string
	StartedAt   time.Time
	FinishesAt  time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
