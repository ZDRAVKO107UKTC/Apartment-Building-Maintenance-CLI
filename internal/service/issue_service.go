package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

type Repository interface {
	Create(issue *model.Issue) error
	FindAll() ([]model.Issue, error)
	FindByID(id uint) (*model.Issue, error)
	Update(issue *model.Issue) error
	Delete(id uint) error
}

type Notifier interface {
	IssueCreated(issue *model.Issue)
	IssueResolved(issue *model.Issue)
}

type IssueService struct {
	repo     Repository
	notifier Notifier
}

func NewIssueService(repo Repository, notifier Notifier) *IssueService {
	return &IssueService{repo: repo, notifier: notifier}
}

var validPriorities = map[string]bool{"low": true, "medium": true, "high": true}

func (s *IssueService) CreateIssue(title, unit, description, priority string) (*model.Issue, error) {
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
	if err := s.repo.Create(issue); err != nil {
		return nil, err
	}
	s.notifier.IssueCreated(issue)
	return issue, nil
}

func (s *IssueService) ListIssues() ([]model.Issue, error) {
	return s.repo.FindAll()
}

func (s *IssueService) ViewIssue(id uint) (*model.Issue, error) {
	if id == 0 {
		return nil, errors.New("id must be a positive integer")
	}
	issue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, model.ErrIssueNotFound) {
			return nil, fmt.Errorf("issue #%d not found", id)
		}
		return nil, err
	}
	return issue, nil
}

func (s *IssueService) UpdateIssue(id uint, status string) (*model.Issue, error) {
	if id == 0 {
		return nil, errors.New("id must be a positive integer")
	}
	newStatus := model.Status(strings.ToLower(strings.TrimSpace(status)))
	switch newStatus {
	case model.StatusOpen, model.StatusInProgress, model.StatusResolved, model.StatusClosed:
	default:
		return nil, errors.New("status must be open, in-progress, resolved, or closed")
	}
	issue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, model.ErrIssueNotFound) {
			return nil, fmt.Errorf("issue #%d not found", id)
		}
		return nil, err
	}
	if issue.Status == newStatus {
		return nil, fmt.Errorf("issue #%d is already %s", id, newStatus)
	}
	if !issue.Status.CanTransitionTo(newStatus) {
		if issue.Status == model.StatusClosed {
			return nil, fmt.Errorf("issue #%d is closed and cannot change status", id)
		}
		return nil, fmt.Errorf("cannot change status from %s to %s", issue.Status, newStatus)
	}
	issue.Status = newStatus
	if err := s.repo.Update(issue); err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *IssueService) ResolveIssue(id uint) (*model.Issue, error) {
	issue, err := s.UpdateIssue(id, string(model.StatusResolved))
	if err != nil {
		return nil, err
	}
	s.notifier.IssueResolved(issue)
	return issue, nil
}

func (s *IssueService) DeleteIssue(id uint) error {
	if id == 0 {
		return errors.New("id must be a positive integer")
	}
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, model.ErrIssueNotFound) {
			return fmt.Errorf("issue #%d not found", id)
		}
		return err
	}
	return s.repo.Delete(id)
}
