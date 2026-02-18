// Package commands provides the CLI commands for Terraship.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the application version
	Version = "dev"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "terraship",
	Short: "Multi-cloud Terraform validation tool",
	Long: `Terraship is a multi-cloud Terraform validation plugin that validates
infrastructure against policies and detects drift.

It supports AWS, Azure, and GCP with two modes:
  - validate-existing: Validates existing infrastructure without applying changes
  - ephemeral-sandbox: Creates temporary infrastructure, validates, and destroys

For more information, visit: https://github.com/terraship/terraship`,
	Version: Version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("Terraship %s (built %s)\n", Version, BuildTime))
}
