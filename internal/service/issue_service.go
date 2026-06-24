package service

import (
	"errors"
	"strings"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/repository"
)

var validPriorities = map[string]bool{"low": true, "medium": true, "high": true}

func CreateIssue(unit, description, priority string) (*model.Issue, error) {
	unit = strings.TrimSpace(unit)
	description = strings.TrimSpace(description)

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
	return issue, nil
}

func ListIssues() ([]model.Issue, error) {
	return repository.FindAllIssues()
}

func ViewIssue(id uint) (*model.Issue, error) {
	return repository.FindIssueByID(id)
}

func UpdateIssue(id uint, status string) (*model.Issue, error) {
	s := model.Status(status)
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
	return issue, nil
}

func DeleteIssue(id uint) error {
	if _, err := repository.FindIssueByID(id); err != nil {
		return err
	}
	return repository.DeleteIssue(id)
}
