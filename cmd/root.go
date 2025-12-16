package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "harden",
	Short: "Policy-driven system hardening tool",
	Long:  "Audit, snapshot, enforce, and rollback system security policies.",
}

func Execute() error {
	return rootCmd.Execute()
}
