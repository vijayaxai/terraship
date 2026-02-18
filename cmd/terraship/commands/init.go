// Package commands provides CLI commands.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new Terraship policy",
	Long: `Initialize a new Terraship policy file in the specified directory.

This command creates a sample policy file that you can customize for your
infrastructure validation needs.

Example:
  terraship init
  terraship init ./my-project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	directory := "."
	if len(args) > 0 {
		directory = args[0]
	}

	policyFile := directory + "/terraship-policy.yml"

	// Check if file already exists
	// For now, just print a message
	fmt.Printf("Initializing Terraship policy in: %s\n", directory)
	fmt.Printf("Policy file will be created at: %s\n", policyFile)
	fmt.Println("\nTo get started:")
	fmt.Println("  1. Copy the sample policy from: policies/sample-policy.yml")
	fmt.Println("  2. Customize the rules for your infrastructure")
	fmt.Println("  3. Run: terraship validate --policy ./terraship-policy.yml")

	return nil
}
