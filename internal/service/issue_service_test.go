package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

// mockRepository is an in-memory stand-in for the data layer. It lets the
// business-logic tests run without a real database or any network call, and can
// be told to return an error to exercise failure paths.
type mockRepository struct {
	issues   map[uint]*model.Issue
	nextID   uint
	failNext error // if set, the next mutating call returns this error
}

func newMockRepository() *mockRepository {
	return &mockRepository{issues: make(map[uint]*model.Issue), nextID: 1}
}

func (m *mockRepository) Create(issue *model.Issue) error {
	if m.failNext != nil {
		err := m.failNext
		m.failNext = nil
		return err
	}
	issue.ID = m.nextID
	m.nextID++
	// Store a copy so callers can't mutate our state through the pointer.
	stored := *issue
	m.issues[issue.ID] = &stored
	return nil
}

func (m *mockRepository) FindAll() ([]model.Issue, error) {
	out := make([]model.Issue, 0, len(m.issues))
	for _, v := range m.issues {
		out = append(out, *v)
	}
	return out, nil
}

func (m *mockRepository) FindByID(id uint) (*model.Issue, error) {
	issue, ok := m.issues[id]
	if !ok {
		return nil, model.ErrIssueNotFound
	}
	copied := *issue
	return &copied, nil
}

func (m *mockRepository) Update(issue *model.Issue) error {
	if m.failNext != nil {
		err := m.failNext
		m.failNext = nil
		return err
	}
	stored := *issue
	m.issues[issue.ID] = &stored
	return nil
}

func (m *mockRepository) Delete(id uint) error {
	delete(m.issues, id)
	return nil
}

// spyNotifier records which notifications the service triggered, so tests can
// assert that email is sent on the right events without touching SendGrid.
type spyNotifier struct {
	created  int
	resolved int
}

func (s *spyNotifier) IssueCreated(*model.Issue)  { s.created++ }
func (s *spyNotifier) IssueResolved(*model.Issue) { s.resolved++ }

// newTestService wires the service to mocks and returns all three so tests can
// make assertions against the collaborators.
func newTestService() (*IssueService, *mockRepository, *spyNotifier) {
	repo := newMockRepository()
	notifier := &spyNotifier{}
	return NewIssueService(repo, notifier), repo, notifier
}

func mustCreate(t *testing.T, svc *IssueService) *model.Issue {
	t.Helper()
	issue, err := svc.CreateIssue("Broken boiler", "4B", "No heat on 4th floor", "high")
	if err != nil {
		t.Fatalf("CreateIssue returned unexpected error: %v", err)
	}
	return issue
}

func TestCreateIssue_Valid(t *testing.T) {
	svc, repo, notifier := newTestService()

	issue, err := svc.CreateIssue("  Broken boiler  ", "  4B ", " No heat ", "high")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.ID == 0 {
		t.Error("expected a non-zero ID after creation")
	}
	if issue.Status != model.StatusOpen {
		t.Errorf("new issue status = %q, want %q", issue.Status, model.StatusOpen)
	}
	if issue.Title != "Broken boiler" || issue.Unit != "4B" || issue.Description != "No heat" {
		t.Errorf("fields not trimmed: %+v", issue)
	}
	if notifier.created != 1 {
		t.Errorf("expected 1 created notification, got %d", notifier.created)
	}
	if len(repo.issues) != 1 {
		t.Errorf("expected issue persisted in repo, got %d", len(repo.issues))
	}
}

func TestCreateIssue_InvalidInput(t *testing.T) {
	cases := []struct {
		name                               string
		title, unit, description, priority string
	}{
		{"empty title", "", "4B", "desc", "high"},
		{"whitespace title", "   ", "4B", "desc", "high"},
		{"empty unit", "Title", "", "desc", "high"},
		{"empty description", "Title", "4B", "", "high"},
		{"invalid priority", "Title", "4B", "desc", "urgent"},
		{"empty priority", "Title", "4B", "desc", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, _, notifier := newTestService()
			issue, err := svc.CreateIssue(tc.title, tc.unit, tc.description, tc.priority)
			if err == nil {
				t.Fatalf("expected an error, got issue %+v", issue)
			}
			if issue != nil {
				t.Errorf("expected nil issue on error, got %+v", issue)
			}
			if notifier.created != 0 {
				t.Errorf("no notification should be sent on validation failure, got %d", notifier.created)
			}
		})
	}
}

func TestCreateIssue_RepositoryError(t *testing.T) {
	svc, repo, notifier := newTestService()
	repo.failNext = errors.New("db is down")

	if _, err := svc.CreateIssue("Title", "4B", "desc", "high"); err == nil {
		t.Fatal("expected the repository error to propagate")
	}
	if notifier.created != 0 {
		t.Error("no notification should be sent when persistence fails")
	}
}

func TestUpdateIssue_Valid(t *testing.T) {
	svc, _, _ := newTestService()
	issue := mustCreate(t, svc)

	updated, err := svc.UpdateIssue(issue.ID, "in-progress")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("status = %q, want %q", updated.Status, model.StatusInProgress)
	}

	reloaded, err := svc.ViewIssue(issue.ID)
	if err != nil {
		t.Fatalf("ViewIssue error: %v", err)
	}
	if reloaded.Status != model.StatusInProgress {
		t.Errorf("persisted status = %q, want %q", reloaded.Status, model.StatusInProgress)
	}
}

func TestUpdateIssue_CaseInsensitiveAndTrimmed(t *testing.T) {
	svc, _, _ := newTestService()
	issue := mustCreate(t, svc)

	updated, err := svc.UpdateIssue(issue.ID, "  IN-PROGRESS  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("status = %q, want %q", updated.Status, model.StatusInProgress)
	}
}

func TestUpdateIssue_StatusValidation(t *testing.T) {
	t.Run("unknown status string", func(t *testing.T) {
		svc, _, _ := newTestService()
		issue := mustCreate(t, svc)
		if _, err := svc.UpdateIssue(issue.ID, "banana"); err == nil {
			t.Fatal("expected error for unknown status")
		}
	})

	t.Run("re-applying current status rejected", func(t *testing.T) {
		svc, _, _ := newTestService()
		issue := mustCreate(t, svc) // starts open
		_, err := svc.UpdateIssue(issue.ID, "open")
		if err == nil {
			t.Fatal("expected error when setting status to its current value")
		}
		if !strings.Contains(err.Error(), "already") {
			t.Errorf("error = %q, want it to mention 'already'", err)
		}
	})

	t.Run("invalid transition from resolved to open", func(t *testing.T) {
		svc, _, _ := newTestService()
		issue := mustCreate(t, svc)
		if _, err := svc.UpdateIssue(issue.ID, "resolved"); err != nil {
			t.Fatalf("setup transition failed: %v", err)
		}
		if _, err := svc.UpdateIssue(issue.ID, "open"); err == nil {
			t.Fatal("expected error transitioning resolved -> open")
		}
	})

	t.Run("closed is terminal", func(t *testing.T) {
		svc, _, _ := newTestService()
		issue := mustCreate(t, svc)
		if _, err := svc.UpdateIssue(issue.ID, "closed"); err != nil {
			t.Fatalf("setup transition to closed failed: %v", err)
		}
		_, err := svc.UpdateIssue(issue.ID, "open")
		if err == nil {
			t.Fatal("expected error changing a closed issue")
		}
		if !strings.Contains(err.Error(), "closed") {
			t.Errorf("error = %q, want it to mention 'closed'", err)
		}
	})
}

func TestUpdateIssue_InvalidInput(t *testing.T) {
	svc, _, _ := newTestService()

	t.Run("zero id", func(t *testing.T) {
		if _, err := svc.UpdateIssue(0, "in-progress"); err == nil {
			t.Fatal("expected error for id 0")
		}
	})

	t.Run("nonexistent id", func(t *testing.T) {
		if _, err := svc.UpdateIssue(9999, "in-progress"); err == nil {
			t.Fatal("expected error for nonexistent id")
		}
	})
}

func TestResolveIssue(t *testing.T) {
	svc, _, notifier := newTestService()
	issue := mustCreate(t, svc)

	resolved, err := svc.ResolveIssue(issue.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.Status != model.StatusResolved {
		t.Errorf("status = %q, want %q", resolved.Status, model.StatusResolved)
	}
	if notifier.resolved != 1 {
		t.Errorf("expected 1 resolved notification, got %d", notifier.resolved)
	}
}

func TestResolveIssue_NotFoundSendsNoNotification(t *testing.T) {
	svc, _, notifier := newTestService()

	if _, err := svc.ResolveIssue(1234); err == nil {
		t.Fatal("expected error resolving nonexistent issue")
	}
	if notifier.resolved != 0 {
		t.Errorf("no notification should be sent when resolve fails, got %d", notifier.resolved)
	}
}

func TestViewIssue_InvalidInput(t *testing.T) {
	svc, _, _ := newTestService()

	t.Run("zero id", func(t *testing.T) {
		if _, err := svc.ViewIssue(0); err == nil {
			t.Fatal("expected error for id 0")
		}
	})

	t.Run("nonexistent id", func(t *testing.T) {
		if _, err := svc.ViewIssue(4242); err == nil {
			t.Fatal("expected error for nonexistent id")
		}
	})
}

func TestDeleteIssue(t *testing.T) {
	t.Run("deletes existing issue", func(t *testing.T) {
		svc, _, _ := newTestService()
		issue := mustCreate(t, svc)
		if err := svc.DeleteIssue(issue.ID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := svc.ViewIssue(issue.ID); err == nil {
			t.Error("expected issue to be gone after delete")
		}
	})

	t.Run("zero id", func(t *testing.T) {
		svc, _, _ := newTestService()
		if err := svc.DeleteIssue(0); err == nil {
			t.Fatal("expected error for id 0")
		}
	})

	t.Run("nonexistent id", func(t *testing.T) {
		svc, _, _ := newTestService()
		if err := svc.DeleteIssue(7777); err == nil {
			t.Fatal("expected error for nonexistent id")
		}
	})
}
