package model

import "testing"

func TestCanTransitionTo(t *testing.T) {
	cases := []struct {
		name   string
		from   Status
		to     Status
		expect bool
	}{
		{"open to in-progress", StatusOpen, StatusInProgress, true},
		{"open to resolved (forward skip)", StatusOpen, StatusResolved, true},
		{"open to closed (forward skip)", StatusOpen, StatusClosed, true},
		{"in-progress back to open", StatusInProgress, StatusOpen, true},
		{"in-progress to resolved", StatusInProgress, StatusResolved, true},
		{"in-progress to closed", StatusInProgress, StatusClosed, true},
		{"resolved reopen to in-progress", StatusResolved, StatusInProgress, true},
		{"resolved to closed", StatusResolved, StatusClosed, true},
		{"resolved back to open not allowed", StatusResolved, StatusOpen, false},
		{"closed is terminal to open", StatusClosed, StatusOpen, false},
		{"closed is terminal to in-progress", StatusClosed, StatusInProgress, false},
		{"closed is terminal to resolved", StatusClosed, StatusResolved, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.from.CanTransitionTo(tc.to); got != tc.expect {
				t.Errorf("%s.CanTransitionTo(%s) = %v, want %v", tc.from, tc.to, got, tc.expect)
			}
		})
	}
}
