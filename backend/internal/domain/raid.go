package domain

import "time"

type Raid struct {
	ID                 string
	AttackerKingdomID  string
	DefenderKingdomID  string
	Status             string
	StartedAt          time.Time
	ArrivesAt          time.Time
	CompletedAt        *time.Time
	Result             *string
	LootJSON           []byte
	AttackerLossesJSON []byte
	DefenderLossesJSON []byte
	ResultJSON         []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type RaidUnit struct {
	ID             string
	RaidID         string
	Side           string
	UnitType       string
	AmountSent     int64
	AmountLost     int64
	AmountReturned int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
