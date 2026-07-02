package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

func newTestNotifier(endpoint string) *SendGridNotifier {
	return &SendGridNotifier{
		client:   &http.Client{Timeout: 2 * time.Second},
		endpoint: endpoint,
		apiKey:   "SG.test-key",
		from:     "no-reply@example.com",
		to:       "manager@example.com",
	}
}

func TestSend_Success(t *testing.T) {
	var gotAuth, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	n := newTestNotifier(srv.URL)
	if err := n.send("subject", "body"); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if gotAuth != "Bearer SG.test-key" {
		t.Errorf("Authorization header = %q, want bearer token", gotAuth)
	}
	if !strings.Contains(gotBody, "manager@example.com") {
		t.Errorf("request body missing recipient: %s", gotBody)
	}
}

func TestSend_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := newTestNotifier(srv.URL)
	err := n.send("subject", "body")
	if err == nil {
		t.Fatal("expected an error for a non-2xx response")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error = %q, want it to mention status 401", err)
	}
}

func TestSend_MissingConfigSkips(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	n := &SendGridNotifier{
		client:   &http.Client{Timeout: 2 * time.Second},
		endpoint: srv.URL,
	}
	if err := n.send("subject", "body"); err != nil {
		t.Fatalf("missing config should skip without error, got: %v", err)
	}
	if called {
		t.Error("no HTTP request should be made when configuration is missing")
	}
}

func TestIssueResolved_DoesNotPanicOnFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	n := newTestNotifier(srv.URL)
	n.IssueResolved(&model.Issue{ID: 1, Title: "x", Unit: "4B"})
}
