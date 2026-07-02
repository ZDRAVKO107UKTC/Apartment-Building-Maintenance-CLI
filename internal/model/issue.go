package model

import (
	"errors"
	"time"
)

var ErrIssueNotFound = errors.New("issue not found")

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in-progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

var validTransitions = map[Status][]Status{
	StatusOpen:       {StatusInProgress, StatusResolved, StatusClosed},
	StatusInProgress: {StatusOpen, StatusResolved, StatusClosed},
	StatusResolved:   {StatusInProgress, StatusClosed},
	StatusClosed:     {},
}

func (s Status) CanTransitionTo(target Status) bool {
	for _, allowed := range validTransitions[s] {
		if allowed == target {
			return true
		}
	}
	return false
}

type Issue struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"not null"`
	Unit        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Priority    string `gorm:"not null"`
	Status      Status `gorm:"not null;default:'open'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
