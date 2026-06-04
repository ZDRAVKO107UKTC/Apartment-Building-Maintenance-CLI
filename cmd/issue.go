package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/database"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/service"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage maintenance issues",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		database.Init()
	},
}

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new maintenance issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		unit, _ := cmd.Flags().GetString("unit")
		description, _ := cmd.Flags().GetString("description")
		priority, _ := cmd.Flags().GetString("priority")

		issue, err := service.CreateIssue(unit, description, priority)
		if err != nil {
			return err
		}
		fmt.Printf("Issue #%d created (unit: %s, priority: %s)\n", issue.ID, issue.Unit, issue.Priority)
		return nil
	},
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all maintenance issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		issues, err := service.ListIssues()
		if err != nil {
			return err
		}
		if len(issues) == 0 {
			fmt.Println("No issues found.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tUNIT\tDESCRIPTION\tPRIORITY\tSTATUS\tCREATED AT")
		for _, issue := range issues {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
				issue.ID, issue.Unit, issue.Description, issue.Priority, issue.Status,
				issue.CreatedAt.Format("2006-01-02 15:04"),
			)
		}
		w.Flush()
		return nil
	},
}

var issueViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a maintenance issue by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetUint("id")

		issue, err := service.ViewIssue(id)
		if err != nil {
			return fmt.Errorf("issue #%d not found", id)
		}
		fmt.Printf("ID:          %d\n", issue.ID)
		fmt.Printf("Unit:        %s\n", issue.Unit)
		fmt.Printf("Description: %s\n", issue.Description)
		fmt.Printf("Priority:    %s\n", issue.Priority)
		fmt.Printf("Status:      %s\n", issue.Status)
		fmt.Printf("Created:     %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:     %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))
		return nil
	},
}

var issueUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the status of a maintenance issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetUint("id")
		status, _ := cmd.Flags().GetString("status")

		issue, err := service.UpdateIssue(id, status)
		if err != nil {
			return err
		}
		fmt.Printf("Issue #%d status updated to: %s\n", issue.ID, issue.Status)
		return nil
	},
}

var issueResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Mark a maintenance issue as resolved",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetUint("id")

		issue, err := service.ResolveIssue(id)
		if err != nil {
			return err
		}
		fmt.Printf("Issue #%d resolved.\n", issue.ID)
		return nil
	},
}

var issueDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a maintenance issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetUint("id")

		if err := service.DeleteIssue(id); err != nil {
			return err
		}
		fmt.Printf("Issue #%d deleted.\n", id)
		return nil
	},
}

func init() {
	issueCreateCmd.Flags().String("unit", "", "Unit identifier (e.g. 4B)")
	issueCreateCmd.Flags().String("description", "", "Description of the issue")
	issueCreateCmd.Flags().String("priority", "", "Priority: low, medium, high")
	_ = issueCreateCmd.MarkFlagRequired("unit")
	_ = issueCreateCmd.MarkFlagRequired("description")
	_ = issueCreateCmd.MarkFlagRequired("priority")

	issueViewCmd.Flags().Uint("id", 0, "Issue ID")
	_ = issueViewCmd.MarkFlagRequired("id")

	issueUpdateCmd.Flags().Uint("id", 0, "Issue ID")
	issueUpdateCmd.Flags().String("status", "", "New status: open, in-progress, resolved")
	_ = issueUpdateCmd.MarkFlagRequired("id")
	_ = issueUpdateCmd.MarkFlagRequired("status")

	issueResolveCmd.Flags().Uint("id", 0, "Issue ID")
	_ = issueResolveCmd.MarkFlagRequired("id")

	issueDeleteCmd.Flags().Uint("id", 0, "Issue ID")
	_ = issueDeleteCmd.MarkFlagRequired("id")

	issueCmd.AddCommand(issueCreateCmd, issueListCmd, issueViewCmd, issueUpdateCmd, issueResolveCmd, issueDeleteCmd)
	rootCmd.AddCommand(issueCmd)
}
