package main

import (
	"fmt"
	"strconv"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/database"
	"github.com/spf13/cobra"
)

// resolveCmd is a top-level shortcut for "issue resolve --id <id>", taking the
// issue ID as a positional argument: `maintenance resolve 1`.
var resolveCmd = &cobra.Command{
	Use:   "resolve <id>",
	Short: "Mark a maintenance issue as resolved (shortcut)",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		database.Init()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid issue ID %q: must be a positive integer", args[0])
		}

		issue, err := newIssueService().ResolveIssue(uint(id))
		if err != nil {
			return err
		}
		fmt.Printf("Issue #%d resolved.\n", issue.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}
