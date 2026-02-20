// Package commands provides CLI commands.
package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijayaxai/terraship/internal/core"
	"github.com/vijayaxai/terraship/internal/output"
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

  # Generate interactive HTML report
  terraship validate ./terraform --output html

  # Generate PDF report (requires wkhtmltopdf)
  terraship validate ./terraform --output pdf

  # Generate all formats
  terraship validate ./terraform --output human,html,pdf,json,sarif

  # Compare with previous run
  terraship validate ./terraform --compare previous-report.json

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
	policyPath     string
	cloudProvider  string
	region         string
	mode           string
	outputFormat   string
	outputFile     string
	noDestroy      bool
	verbose        bool
	htmlAdvanced   bool
	includeHistory bool
	compareWith    string
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&policyPath, "policy", "p", "./policies/sample-policy.yml", "Path to policy YAML file")
	validateCmd.Flags().StringVar(&cloudProvider, "provider", "", "Cloud provider (aws, azure, gcp) - auto-detect if not specified")
	validateCmd.Flags().StringVar(&region, "region", "", "Cloud region (AWS region, Azure location, GCP region)")
	validateCmd.Flags().StringVarP(&mode, "mode", "m", "validate-existing", "Validation mode: validate-existing or ephemeral-sandbox")
	validateCmd.Flags().StringVarP(&outputFormat, "output", "o", "human", "Output format: human, json, html, pdf, sarif (comma-separated for multiple)")
	validateCmd.Flags().StringVarP(&outputFile, "output-file", "f", "", "Write output to file instead of stdout")
	validateCmd.Flags().BoolVar(&noDestroy, "no-destroy", false, "Don't destroy resources in ephemeral mode")
	validateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	validateCmd.Flags().BoolVar(&htmlAdvanced, "html-advanced", false, "Use advanced HTML features (dark mode, charts, search)")
	validateCmd.Flags().BoolVar(&includeHistory, "include-history", false, "Include validation history in report")
	validateCmd.Flags().StringVar(&compareWith, "compare", "", "Compare with previous validation results (JSON file)")
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
		return fmt.Errorf("policy file does not exist: %s\n\n"+
			"Create a policy file by running:\n"+
			"  terraship init\n\n"+
			"Or specify a custom policy with:\n"+
			"  terraship validate . --policy ./your-policy.yml\n\n"+
			"For help with policy files, see:\n"+
			"  terraship validate --help", policyPath)
	}

	// Validate mode
	if mode != "validate-existing" && mode != "ephemeral-sandbox" {
		return fmt.Errorf("invalid mode: %s (must be validate-existing or ephemeral-sandbox)", mode)
	}

	// Validate output formats
	formats := strings.Split(outputFormat, ",")
	for _, f := range formats {
		f = strings.TrimSpace(f)
		if f != "human" && f != "json" && f != "html" && f != "pdf" && f != "sarif" {
			return fmt.Errorf("invalid output format: %s (must be human, json, html, pdf, or sarif)", f)
		}
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

	// Convert summary to ValidationResult for report generation
	validationResult := convertSummaryToValidationResult(summary)

	// Load previous results if comparing
	var previousResults *output.ValidationResult
	if compareWith != "" {
		prevResults, err := loadValidationResultsFromFile(compareWith)
		if err != nil {
			fmt.Printf("⚠  Warning: Could not load previous results: %v\n", err)
		} else {
			previousResults = prevResults
		}
	}

	// Process each output format
	for _, f := range formats {
		f = strings.TrimSpace(f)
		if err := generateValidationReport(f, validationResult, previousResults); err != nil {
			fmt.Printf("❌ Error generating %s report: %v\n", f, err)
			continue
		}
	}

	// Print summary if not only outputting to a file
	if outputFile == "" || strings.Contains(outputFormat, "human") {
		printValidationSummary(validationResult)
	}

	// Exit with error code if validation failed
	if summary.FailedResources > 0 || summary.ErrorResources > 0 {
		os.Exit(1)
	}

	return nil
}

// convertSummaryToValidationResult converts core validator summary to ValidationResult
func convertSummaryToValidationResult(summary *core.Summary) *output.ValidationResult {
	result := &output.ValidationResult{
		TotalResources:   summary.TotalResources,
		PassedResources:  summary.PassedResources,
		FailedResources:  summary.FailedResources,
		WarningResources: summary.WarningResources,
		Timestamp:        time.Now().Format("2006-01-02 15:04:05"),
		Resources:        convertResourcesToOutputFormat(summary),
	}
	return result
}

// convertResourcesToOutputFormat converts core resources to output resources
func convertResourcesToOutputFormat(summary *core.Summary) []output.Resource {
	resources := make([]output.Resource, 0)
	
	for _, report := range summary.Reports {
		// Create resource
		resource := output.Resource{
			Name:        report.ResourceAddress,
			Type:        report.ResourceType,
			Provider:    report.Provider,
			IsFailed:    report.Status == "fail" || report.Status == "error",
			HasWarnings: report.Status == "warning",
		}
		
		// Convert rule results to checks
		for _, result := range report.RuleResults {
			check := output.Check{
				Name:        result.RuleName,
				Message:     result.Message,
				Severity:    result.Severity,
				Failed:      !result.Passed,
				Warning:     result.Severity == "warning" && result.Passed,
				Remediation: result.Remediation,
				Details:     result.Details,
			}
			
			resource.Checks = append(resource.Checks, check)
		}
		
		// Add errors as checks if any
		for _, errMsg := range report.Errors {
			check := output.Check{
				Name:     "Validation Error",
				Message:  errMsg,
				Severity: "error",
				Failed:   true,
			}
			resource.Checks = append(resource.Checks, check)
		}
		
		resources = append(resources, resource)
	}
	
	return resources
}

// generateValidationReport generates report in specified format
func generateValidationReport(format string, results *output.ValidationResult, previousResults *output.ValidationResult) error {
	switch format {
	case "html":
		return generateHTMLReport(results, previousResults)
	case "pdf":
		return generatePDFReport(results, previousResults)
	case "json":
		return generateJSONReportFile(results)
	case "sarif":
		return generateSARIFReportFile(results)
	case "human":
		printHumanReport(results)
		return nil
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

// generateHTMLReport creates HTML report
func generateHTMLReport(results *output.ValidationResult, previousResults *output.ValidationResult) error {
	// Generate HTML
	html, err := output.GenerateHTML(results, includeHistory, previousResults)
	if err != nil {
		fmt.Printf("❌ Error generating html report: %v\n", err)
		return err
	}

	// Determine output file
	outFile := outputFile
	if outFile == "" {
		outFile = "report.html"
	}

	// Save HTML to file
	if err := os.WriteFile(outFile, []byte(html), 0644); err != nil {
		return err
	}

	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	fmt.Printf("%s✓%s HTML report generated: %s\n", colorGreen, colorReset, outFile)
	fmt.Printf("  Open with: open %s (macOS) or your web browser\n", outFile)

	return nil
}

// generatePDFReport creates PDF report
func generatePDFReport(results *output.ValidationResult, previousResults *output.ValidationResult) error {
	// For now, generate HTML and inform user to export as PDF from browser
	html, err := output.GenerateHTML(results, includeHistory, previousResults)
	if err != nil {
		return err
	}

	// Save HTML with .pdf extension suggestion
	outFile := outputFile
	if outFile == "" {
		outFile = "report.html"
	}

	if err := os.WriteFile(outFile, []byte(html), 0644); err != nil {
		return err
	}

	colorYellow := "\033[93m"
	colorReset := "\033[0m"
	fmt.Printf("%s⚠%s PDF export requires wkhtmltopdf or browser export\n", colorYellow, colorReset)
	fmt.Printf("  HTML report saved: %s\n", outFile)
	fmt.Printf("  To convert to PDF:\n")
	fmt.Printf("    1. Open in browser: open %s\n", outFile)
	fmt.Printf("    2. Press Ctrl+P (or Cmd+P) and save as PDF\n")
	fmt.Printf("    OR install wkhtmltopdf: brew install wkhtmltopdf (macOS)\n")

	return nil
}

// generateJSONReportFile creates JSON report file
func generateJSONReportFile(results *output.ValidationResult) error {
	outFile := outputFile
	if outFile == "" {
		outFile = "terraship-report.json"
	}

	jsonBytes, err := results.ToJSON()
	if err != nil {
		return err
	}

	if err := os.WriteFile(outFile, jsonBytes, 0644); err != nil {
		return err
	}

	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	fmt.Printf("%s✓%s JSON report generated: %s\n", colorGreen, colorReset, outFile)

	return nil
}

// generateSARIFReportFile creates SARIF report file
func generateSARIFReportFile(results *output.ValidationResult) error {
	outFile := outputFile
	if outFile == "" {
		outFile = "terraship-report.sarif"
	}

	sarifBytes, err := results.ToSARIF()
	if err != nil {
		return err
	}

	if err := os.WriteFile(outFile, sarifBytes, 0644); err != nil {
		return err
	}

	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	fmt.Printf("%s✓%s SARIF report generated: %s\n", colorGreen, colorReset, outFile)

	return nil
}

// printHumanReport prints human-readable report
func printHumanReport(results *output.ValidationResult) {
	fmt.Println("\n" + strings.Repeat("=", 63))
	fmt.Println("                    TERRASHIP VALIDATION REPORT")
	fmt.Println(strings.Repeat("=", 63))
	fmt.Println()
	fmt.Printf("SUMMARY:\n")
	fmt.Printf("  Total Resources:    %d\n", results.TotalResources)
	fmt.Printf("  ✓ Passed:           %d\n", results.PassedResources)
	fmt.Printf("  ✗ Failed:           %d\n", results.FailedResources)
	fmt.Printf("  ⚠ Warnings:         %d\n", results.WarningResources)
	fmt.Println()

	if results.FailedResources > 0 {
		fmt.Println("✗ VALIDATION FAILED")
	} else {
		fmt.Println("✓ VALIDATION PASSED")
	}
}

// printValidationSummary prints summary statistics
func printValidationSummary(results *output.ValidationResult) {
	compliance := 0.0
	if results.TotalResources > 0 {
		compliance = (float64(results.PassedResources) / float64(results.TotalResources)) * 100
	}

	fmt.Printf("\n📊 Compliance Score: %.1f%%\n", compliance)
	fmt.Printf("⏱  Validation completed: %s\n\n", results.Timestamp)
}

// loadValidationResultsFromFile loads previous validation results
func loadValidationResultsFromFile(filePath string) (*output.ValidationResult, error) {
	// This would load and parse the previous results file
	// For now, placeholder
	return nil, fmt.Errorf("loading previous results not yet implemented")
}
