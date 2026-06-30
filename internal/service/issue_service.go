package service

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/repository"
)

var validPriorities = map[string]bool{"low": true, "medium": true, "high": true}

func CreateIssue(title, unit, description, priority string) (*model.Issue, error) {
	title = strings.TrimSpace(title)
	unit = strings.TrimSpace(unit)
	description = strings.TrimSpace(description)

	if title == "" {
		return nil, errors.New("title must not be empty")
	}
	if unit == "" {
		return nil, errors.New("unit must not be empty")
	}
	if description == "" {
		return nil, errors.New("description must not be empty")
	}
	if !validPriorities[priority] {
		return nil, errors.New("priority must be low, medium, or high")
	}
	issue := &model.Issue{
		Title:       title,
		Unit:        unit,
		Description: description,
		Priority:    priority,
		Status:      model.StatusOpen,
	}
	if err := repository.CreateIssue(issue); err != nil {
		return nil, err
	}
	NotifyIssueCreated(issue)
	return issue, nil
}

func ListIssues() ([]model.Issue, error) {
	return repository.FindAllIssues()
}

func ViewIssue(id uint) (*model.Issue, error) {
	if id == 0 {
		return nil, errors.New("id must be a positive integer")
	}
	issue, err := repository.FindIssueByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("issue #%d not found", id)
		}
		return nil, err
	}
	return issue, nil
}

func UpdateIssue(id uint, status string) (*model.Issue, error) {
	if id == 0 {
		return nil, errors.New("id must be a positive integer")
	}
	s := model.Status(strings.ToLower(strings.TrimSpace(status)))
	switch s {
	case model.StatusOpen, model.StatusInProgress, model.StatusResolved, model.StatusClosed:
	default:
		return nil, errors.New("status must be open, in-progress, resolved, or closed")
	}
	issue, err := repository.FindIssueByID(id)
	if err != nil {
		return nil, err
	}
	issue.Status = s
	if err := repository.UpdateIssue(issue); err != nil {
		return nil, err
	}
	return issue, nil
}

func ResolveIssue(id uint) (*model.Issue, error) {
	issue, err := repository.FindIssueByID(id)
	if err != nil {
		return nil, err
	}
	issue.Status = model.StatusResolved
	if err := repository.UpdateIssue(issue); err != nil {
		return nil, err
	}
	NotifyIssueResolved(issue)
	return issue, nil
}

func DeleteIssue(id uint) error {
	if _, err := repository.FindIssueByID(id); err != nil {
		return err
	}
	return repository.DeleteIssue(id)
}
