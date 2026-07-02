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

const sendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"

type SendGridNotifier struct {
	client   *http.Client
	endpoint string
	apiKey   string
	from     string
	to       string
}

func NewSendGridNotifier() *SendGridNotifier {
	return &SendGridNotifier{
		client:   &http.Client{Timeout: 10 * time.Second},
		endpoint: sendGridEndpoint,
		apiKey:   os.Getenv("SENDGRID_API_KEY"),
		from:     os.Getenv("EMAIL_FROM"),
		to:       os.Getenv("MANAGER_EMAIL"),
	}
}

func (n *SendGridNotifier) IssueCreated(issue *model.Issue) {
	subject := fmt.Sprintf("New maintenance issue #%d: %s", issue.ID, issue.Title)
	body := fmt.Sprintf(
		"A new maintenance issue has been registered.\n\n"+
			"ID:          %d\nTitle:       %s\nUnit:        %s\nPriority:    %s\nStatus:      %s\nDescription: %s\n",
		issue.ID, issue.Title, issue.Unit, issue.Priority, issue.Status, issue.Description,
	)
	if err := n.send(subject, body); err != nil {
		log.Printf("email notification failed: %v", err)
	}
}

func (n *SendGridNotifier) IssueResolved(issue *model.Issue) {
	subject := fmt.Sprintf("Maintenance issue #%d resolved: %s", issue.ID, issue.Title)
	body := fmt.Sprintf(
		"The following maintenance issue has been marked resolved.\n\n"+
			"ID:    %d\nTitle: %s\nUnit:  %s\n",
		issue.ID, issue.Title, issue.Unit,
	)
	if err := n.send(subject, body); err != nil {
		log.Printf("email notification failed: %v", err)
	}
}

func (n *SendGridNotifier) send(subject, body string) error {
	if n.apiKey == "" || n.from == "" || n.to == "" {
		log.Println("email notification skipped: SENDGRID_API_KEY, EMAIL_FROM or MANAGER_EMAIL not configured")
		return nil
	}

	payload := map[string]any{
		"personalizations": []map[string]any{
			{"to": []map[string]string{{"email": n.to}}},
		},
		"from":    map[string]string{"email": n.from},
		"subject": subject,
		"content": []map[string]string{
			{"type": "text/plain", "value": body},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not encode payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), n.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("SendGrid returned status %d", resp.StatusCode)
	}
	return nil
}
