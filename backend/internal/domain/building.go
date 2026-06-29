package domain

import "time"

type Building struct {
	ID                string
	KingdomID         string
	Type              string
	Level             int
	UpgradeStartedAt  *time.Time
	UpgradeFinishesAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (b Building) IsUpgrading() bool {
	return b.UpgradeStartedAt != nil && b.UpgradeFinishesAt != nil
}
