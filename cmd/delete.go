package main

import (
	"fmt"
	"strconv"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/database"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/service"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a maintenance issue (shortcut)",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		database.Init()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid issue ID %q: must be a positive integer", args[0])
		}

		if err := service.DeleteIssue(uint(id)); err != nil {
			return err
		}
		fmt.Printf("Issue #%d deleted.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
