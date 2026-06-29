package domain

import "time"

type MissionReport struct {
	ID          string
	KingdomID   string
	MissionID   *string
	Type        string
	Title       string
	Body        string
	Result      string
	RewardsJSON []byte
	LossesJSON  []byte
	IsRead      bool
	CreatedAt   time.Time
}
