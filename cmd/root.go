package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/database"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/repository"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/service"
)

var rootCmd = &cobra.Command{
	Use:   "maintenance",
	Short: "Apartment Building Maintenance CLI",
	Long:  "A CLI tool for managing apartment building maintenance issues.",
	// Runtime errors are reported once by Execute (below); silence Cobra's own
	// error/usage dump so failures print a single meaningful line to stderr.
	SilenceUsage:  true,
	SilenceErrors: true,
}

// newIssueService assembles the service with its real dependencies. It must be
// called after database.Init() has populated database.DB.
func newIssueService() *service.IssueService {
	return service.NewIssueService(
		repository.NewGormIssueRepository(database.DB),
		service.NewSendGridNotifier(),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
