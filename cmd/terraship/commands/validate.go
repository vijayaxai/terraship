// Package commands provides CLI commands.
package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terraship/terraship/internal/core"
	"github.com/terraship/terraship/internal/output"
)

var validateCmd = &cobra.Command{
	Use:   "validate [directory]",
	Short: "Validate Terraform infrastructure",
	Long: `Validate Terraform infrastructure against policy rules.

This command runs Terraform validation and checks resources against
the specified policy file. It supports two modes:

  validate-existing: Validate existing infrastructure without making changes
  ephemeral-sandbox: Create temporary infrastructure, validate, and destroy

Examples:
  # Validate with default policy
  terraship validate ./terraform

  # Validate existing infrastructure
  terraship validate ./terraform --mode validate-existing

  # Create ephemeral environment for testing
  terraship validate ./terraform --mode ephemeral-sandbox

  # Use custom policy and output format
  terraship validate ./terraform --policy ./my-policy.yml --output json

  # Manually specify cloud provider
  terraship validate ./terraform --provider aws --region us-west-2`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

var (
	policyPath    string
	cloudProvider string
	region        string
	mode          string
	outputFormat  string
	outputFile    string
	noDestroy     bool
	verbose       bool
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&policyPath, "policy", "p", "./policies/sample-policy.yml", "Path to policy YAML file")
	validateCmd.Flags().StringVar(&cloudProvider, "provider", "", "Cloud provider (aws, azure, gcp) - auto-detect if not specified")
	validateCmd.Flags().StringVar(&region, "region", "", "Cloud region (AWS region, Azure location, GCP region)")
	validateCmd.Flags().StringVarP(&mode, "mode", "m", "validate-existing", "Validation mode: validate-existing or ephemeral-sandbox")
	validateCmd.Flags().StringVarP(&outputFormat, "output", "o", "human", "Output format: human, json, sarif")
	validateCmd.Flags().StringVarP(&outputFile, "output-file", "f", "", "Write output to file instead of stdout")
	validateCmd.Flags().BoolVar(&noDestroy, "no-destroy", false, "Don't destroy resources in ephemeral mode")
	validateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func runValidate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get working directory
	workingDir := "."
	if len(args) > 0 {
		workingDir = args[0]
	}

	// Validate working directory exists
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", workingDir)
	}

	// Validate policy file exists
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		return fmt.Errorf("policy file does not exist: %s", policyPath)
	}

	// Validate mode
	if mode != "validate-existing" && mode != "ephemeral-sandbox" {
		return fmt.Errorf("invalid mode: %s (must be validate-existing or ephemeral-sandbox)", mode)
	}

	// Validate output format
	if outputFormat != "human" && outputFormat != "json" && outputFormat != "sarif" {
		return fmt.Errorf("invalid output format: %s (must be human, json, or sarif)", outputFormat)
	}

	if verbose {
		fmt.Printf("Starting Terraship validation...\n")
		fmt.Printf("  Working directory: %s\n", workingDir)
		fmt.Printf("  Policy file: %s\n", policyPath)
		fmt.Printf("  Mode: %s\n", mode)
		fmt.Printf("  Output format: %s\n", outputFormat)
		if cloudProvider != "" {
			fmt.Printf("  Cloud provider: %s\n", cloudProvider)
		}
		fmt.Println()
	}

	// Create validator config
	config := core.ValidatorConfig{
		Mode:          core.ValidationMode(mode),
		WorkingDir:    workingDir,
		PolicyPath:    policyPath,
		CloudProvider: cloudProvider,
		OutputFormat:  outputFormat,
		OutputFile:    outputFile,
		NoDestroy:     noDestroy,
		Verbose:       verbose,
	}

	// Create validator
	validator, err := core.NewValidator(config)
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}

	// Run validation
	summary, err := validator.Validate(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Format and output results
	var formatter output.Formatter
	switch outputFormat {
	case "json":
		formatter = output.NewJSONFormatter(true)
	case "sarif":
		formatter = output.NewSARIFFormatter()
	default:
		formatter = output.NewHumanFormatter()
	}

	result, err := formatter.Format(summary)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if verbose {
			fmt.Printf("Results written to: %s\n", outputFile)
		}
	} else {
		fmt.Print(result)
	}

	// Exit with error code if validation failed
	if summary.FailedResources > 0 || summary.ErrorResources > 0 {
		os.Exit(1)
	}

	return nil
}
