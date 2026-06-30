package domain

import "time"

type GameEvent struct {
	ID                  string
	Key                 string
	Category            string
	Title               string
	Body                string
	TriggerType         string
	Weight              int
	IsActive            bool
	CooldownSeconds     int
	ExpiresAfterSeconds int
	ConditionsJSON      []byte
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type EventChoice struct {
	ID          string
	GameEventID string
	Key         string
	Label       string
	Description string
	EffectsJSON []byte
	ResultTitle string
	ResultBody  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type KingdomEvent struct {
	ID                string
	KingdomID         string
	GameEventID       string
	Status            string
	GeneratedAt       time.Time
	ExpiresAt         time.Time
	ResolvedAt        *time.Time
	SelectedChoiceKey *string
	ResultJSON        []byte
	CreatedAt         time.Time
	UpdatedAt         time.Time
	GameEvent         GameEvent
	Choices           []EventChoice
}
