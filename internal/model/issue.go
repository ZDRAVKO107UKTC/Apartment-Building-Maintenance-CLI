package model

import "time"

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in-progress"
	StatusResolved   Status = "resolved"
)

type Issue struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Unit        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Priority    string `gorm:"not null"`
	Status      Status `gorm:"not null;default:'open'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
