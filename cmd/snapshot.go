package cmd

import (
	"fmt"

	"harden/policy"
	"harden/utils"
	"harden/utils/ensure"

	"github.com/spf13/cobra"
)

var snapshotPolicyPath string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create a system snapshot for rollback",
	Long:  "Collects current system state based on policies and stores a rollback snapshot.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSnapshot()
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)

	snapshotCmd.Flags().StringVarP(
		&snapshotPolicyPath,
		"policy",
		"p",
		"./policies/linux.yaml",
		"Path to policy YAML file",
	)
}

func runSnapshot() error {
	// Detect OS
	res, err := utils.Detect()
	if err != nil {
		return err
	}

	fmt.Printf("Detected: %s\n", res.PrettyName)
	fmt.Printf("ID: %s\n", res.ID)
	fmt.Printf("Family: %s\n", res.Family)
	fmt.Printf("Version ID: %s\n", res.VersionID)

	// Enforce privilege here (NOT globally)
	ensure.RequireAdmin()

	// Load policies
	ps, err := policy.LoadFromFile(snapshotPolicyPath)
	if err != nil {
		return err
	}

	// Build snapshot
	snap, err := utils.BuildSnapshot(ps.Policies)
	if err != nil {
		return err
	}

	// Save snapshot
	path, err := snap.SaveNextToExecutable()
	if err != nil {
		return err
	}

	fmt.Println("Snapshot saved at:", path)
	return nil
}
