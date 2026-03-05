package model

import "time"

type Project struct {
	ID        int
	Name      string
	OwnerID   int
	CreatedAt time.Time
}
