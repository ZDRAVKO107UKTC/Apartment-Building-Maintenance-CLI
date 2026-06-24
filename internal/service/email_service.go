package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

// sendGridEndpoint is the SendGrid v3 transactional mail API.
const sendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"

// emailClient is the resilient HTTP client used for all outbound notifications.
// The timeout guarantees a slow or unresponsive SendGrid never blocks the CLI.
var emailClient = &http.Client{Timeout: 10 * time.Second}

// NotifyIssueCreated sends a notification that a new issue was registered.
// Failures are logged but never returned, so notification problems can never
// crash the application or roll back a successful database write.
func NotifyIssueCreated(issue *model.Issue) {
	subject := fmt.Sprintf("New maintenance issue #%d: %s", issue.ID, issue.Title)
	body := fmt.Sprintf(
		"A new maintenance issue has been registered.\n\n"+
			"ID:          %d\nTitle:       %s\nUnit:        %s\nPriority:    %s\nStatus:      %s\nDescription: %s\n",
		issue.ID, issue.Title, issue.Unit, issue.Priority, issue.Status, issue.Description,
	)
	sendEmail(subject, body)
}

// NotifyIssueResolved sends a notification that an issue has been resolved.
// As with NotifyIssueCreated, failures are logged and swallowed.
func NotifyIssueResolved(issue *model.Issue) {
	subject := fmt.Sprintf("Maintenance issue #%d resolved: %s", issue.ID, issue.Title)
	body := fmt.Sprintf(
		"The following maintenance issue has been marked resolved.\n\n"+
			"ID:    %d\nTitle: %s\nUnit:  %s\n",
		issue.ID, issue.Title, issue.Unit,
	)
	sendEmail(subject, body)
}

// sendEmail posts a transactional message to SendGrid. It is deliberately
// best-effort: every failure path logs to stderr and returns without panicking
// or propagating an error to the caller.
func sendEmail(subject, body string) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("MANAGER_EMAIL")

	// In local/dev environments the integration is optional. Skip quietly with
	// a notice rather than treating a missing key as an error.
	if apiKey == "" || from == "" || to == "" {
		log.Println("email notification skipped: SENDGRID_API_KEY, EMAIL_FROM or MANAGER_EMAIL not configured")
		return
	}

	payload := map[string]any{
		"personalizations": []map[string]any{
			{"to": []map[string]string{{"email": to}}},
		},
		"from":    map[string]string{"email": from},
		"subject": subject,
		"content": []map[string]string{
			{"type": "text/plain", "value": body},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("email notification failed: could not encode payload: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), emailClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendGridEndpoint, bytes.NewReader(data))
	if err != nil {
		log.Printf("email notification failed: could not build request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := emailClient.Do(req)
	if err != nil {
		// Covers timeouts, DNS errors, connection refused, etc.
		log.Printf("email notification failed: request error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("email notification failed: SendGrid returned status %d", resp.StatusCode)
		return
	}
}
