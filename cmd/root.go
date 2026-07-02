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
	Use:           "maintenance",
	Short:         "Apartment Building Maintenance CLI",
	Long:          "A CLI tool for managing apartment building maintenance issues.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

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
