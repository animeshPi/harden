package cmd

import (
	"fmt"
	"os"

	"harden/policy"
	"harden/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var auditPolicyPath string

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit system against security policies",
	Long:  "Runs audit checks defined in policy files without applying remediation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAudit()
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	auditCmd.Flags().StringVarP(
		&auditPolicyPath,
		"policy",
		"p",
		"./policies/linux.yaml",
		"Path to policy YAML file",
	)
}

func runAudit() error {
	ps, err := policy.LoadFromFile(auditPolicyPath)
	if err != nil {
		return fmt.Errorf("policy load failed: %w", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{
		"ID",
		"TITLE",
		"SEVERITY",
		"STATUS",
	})

	hasFailures := false

	for _, p := range ps.Policies {
		status := "PASS"

		res := utils.Execute(p.CheckCmd)
		if !res.Success {
			status = "FAIL"
			hasFailures = true
		}

		t.AppendRow(table.Row{
			p.ID,
			p.Title,
			p.Severity,
			status,
		})
	}

	t.Render()

	if hasFailures {
		os.Exit(2)
	}

	return nil
}
