package domain

import "time"

type Ruler struct {
	ID           string
	KingdomID    string
	Name         string
	Age          int
	Culture      string
	Authority    int
	Courage      int
	Cunning      int
	Honor        int
	Cruelty      int
	Ambition     int
	Paranoia     int
	HealthStatus string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
