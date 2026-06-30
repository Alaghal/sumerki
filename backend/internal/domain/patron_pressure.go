package domain

import "time"

type PatronPressureState struct {
	ID                   string
	KingdomID            string
	Patron               string
	TributeDebtGold      int64
	TributeDebtFood      int64
	ContributionDebtFood int64
	PressureLevel        int
	CrisisStatus         string
	CrisisStartedAt      *time.Time
	NextTributeAt        time.Time
	LastResolvedAt       time.Time
	DelayUntil           *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
