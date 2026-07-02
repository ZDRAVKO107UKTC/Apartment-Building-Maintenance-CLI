package model

import "time"

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in-progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

// validTransitions defines which target statuses each status may move to. It
// follows the RFC workflow (open → in-progress → resolved → closed) while
// allowing forward skips and reopening a resolved issue; closed is terminal.
var validTransitions = map[Status][]Status{
	StatusOpen:       {StatusInProgress, StatusResolved, StatusClosed},
	StatusInProgress: {StatusOpen, StatusResolved, StatusClosed},
	StatusResolved:   {StatusInProgress, StatusClosed},
	StatusClosed:     {},
}

// CanTransitionTo reports whether an issue may move from its current status to
// the target status under the workflow rules.
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
