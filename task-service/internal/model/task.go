package model

import "time"

type Task struct {
	ID          int
	ProjectID   int
	Title       string
	Description string
	Status      string
	AssigneeID  int
	CreatedAt   time.Time
}
