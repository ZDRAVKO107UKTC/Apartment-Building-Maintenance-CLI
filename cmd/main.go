package main

import (
	"fmt"
	"os"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "maintenance",
	Short: "Apartment Building Maintenance CLI",
	Long:  "A CLI tool for managing apartment building maintenance issues.",
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: .env file not found")
	}

	database.Init()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
