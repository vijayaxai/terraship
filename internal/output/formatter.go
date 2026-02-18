// Package output provides formatters for validation results.
package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/vijayaxai/terraship/internal/core"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	Format(summary *core.Summary) (string, error)
}

// HumanFormatter produces human-readable output
type HumanFormatter struct{}

// NewHumanFormatter creates a new human formatter
func NewHumanFormatter() *HumanFormatter {
	return &HumanFormatter{}
}

// Format generates human-readable output
func (f *HumanFormatter) Format(summary *core.Summary) (string, error) {
	var sb strings.Builder

	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString("                    TERRASHIP VALIDATION REPORT                  \n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	// Summary section
	sb.WriteString("SUMMARY:\n")
	sb.WriteString(fmt.Sprintf("  Total Resources:    %d\n", summary.TotalResources))
	sb.WriteString(fmt.Sprintf("  ✓ Passed:           %d\n", summary.PassedResources))
	sb.WriteString(fmt.Sprintf("  ✗ Failed:           %d\n", summary.FailedResources))
	sb.WriteString(fmt.Sprintf("  ⚠ Warnings:         %d\n", summary.WarningResources))
	sb.WriteString(fmt.Sprintf("  ⨯ Errors:           %d\n", summary.ErrorResources))
	sb.WriteString(fmt.Sprintf("  ↔ Drift Detected:   %d\n\n", summary.DriftDetected))

	// Overall status
	if summary.FailedResources == 0 && summary.ErrorResources == 0 {
		sb.WriteString("✓ VALIDATION PASSED\n\n")
	} else {
		sb.WriteString("✗ VALIDATION FAILED\n\n")
	}

	// Detailed results
	if len(summary.Reports) > 0 {
		sb.WriteString("DETAILED RESULTS:\n")
		sb.WriteString("───────────────────────────────────────────────────────────────\n\n")

		for _, report := range summary.Reports {
			statusIcon := "✓"
			switch report.Status {
			case "fail":
				statusIcon = "✗"
			case "warning":
				statusIcon = "⚠"
			case "error":
				statusIcon = "⨯"
			}

			sb.WriteString(fmt.Sprintf("%s %s (%s)\n", statusIcon, report.ResourceAddress, report.ResourceType))
			sb.WriteString(fmt.Sprintf("  Provider: %s\n", report.Provider))

			// Rule results
			if len(report.RuleResults) > 0 {
				sb.WriteString("  Policy Checks:\n")
				for _, result := range report.RuleResults {
					resultIcon := "✓"
					if !result.Passed {
						resultIcon = "✗"
					}
					sb.WriteString(fmt.Sprintf("    %s %s [%s]\n", resultIcon, result.RuleName, result.Severity))
					if !result.Passed {
						sb.WriteString(fmt.Sprintf("      Message: %s\n", result.Message))
						for _, detail := range result.Details {
							sb.WriteString(fmt.Sprintf("      - %s\n", detail))
						}
						if result.Remediation != "" {
							sb.WriteString(fmt.Sprintf("      💡 Remediation: %s\n", result.Remediation))
						}
					}
				}
			}

			// Drift detection
			if report.DriftStatus != nil && report.DriftStatus.DriftDetected {
				sb.WriteString("  ↔ Drift Detected:\n")
				for _, detail := range report.DriftStatus.DriftDetails {
					sb.WriteString(fmt.Sprintf("    - %s\n", detail))
				}
			}

			// Errors
			if len(report.Errors) > 0 {
				sb.WriteString("  Errors:\n")
				for _, err := range report.Errors {
					sb.WriteString(fmt.Sprintf("    - %s\n", err))
				}
			}

			sb.WriteString("\n")
		}
	}

	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339)))

	return sb.String(), nil
}

// JSONFormatter produces JSON output
type JSONFormatter struct {
	Pretty bool
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{Pretty: pretty}
}

// Format generates JSON output
func (f *JSONFormatter) Format(summary *core.Summary) (string, error) {
	var data []byte
	var err error

	if f.Pretty {
		data, err = json.MarshalIndent(summary, "", "  ")
	} else {
		data, err = json.Marshal(summary)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(data), nil
}

// SARIFFormatter produces SARIF (Static Analysis Results Interchange Format) output
type SARIFFormatter struct{}

// NewSARIFFormatter creates a new SARIF formatter
func NewSARIFFormatter() *SARIFFormatter {
	return &SARIFFormatter{}
}

// SARIFReport represents a SARIF 2.1.0 report
type SARIFReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SARIFRun `json:"runs"`
}

// SARIFRun represents a SARIF run
type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

// SARIFTool represents the tool information
type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

// SARIFDriver represents the tool driver
type SARIFDriver struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	InformationURI string `json:"informationUri"`
}

// SARIFResult represents a single result
type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"` // "error", "warning", "note"
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
}

// SARIFMessage represents a result message
type SARIFMessage struct {
	Text string `json:"text"`
}

// SARIFLocation represents a result location
type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

// SARIFPhysicalLocation represents physical location
type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
}

// SARIFArtifactLocation represents an artifact location
type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

// Format generates SARIF output
func (f *SARIFFormatter) Format(summary *core.Summary) (string, error) {
	sarif := SARIFReport{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "Terraship",
						Version:        "1.0.0",
						InformationURI: "https://github.com/vijayaxai/terraship",
					},
				},
				Results: []SARIFResult{},
			},
		},
	}

	// Convert validation results to SARIF results
	for _, report := range summary.Reports {
		for _, result := range report.RuleResults {
			if !result.Passed {
				level := "warning"
				if result.Severity == "error" {
					level = "error"
				} else if result.Severity == "info" {
					level = "note"
				}

				message := result.Message
				if len(result.Details) > 0 {
					message += "\n" + strings.Join(result.Details, "\n")
				}

				sarifResult := SARIFResult{
					RuleID: result.RuleName,
					Level:  level,
					Message: SARIFMessage{
						Text: message,
					},
					Locations: []SARIFLocation{
						{
							PhysicalLocation: SARIFPhysicalLocation{
								ArtifactLocation: SARIFArtifactLocation{
									URI: report.ResourceAddress,
								},
							},
						},
					},
				}

				sarif.Runs[0].Results = append(sarif.Runs[0].Results, sarifResult)
			}
		}
	}

	data, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal SARIF: %w", err)
	}

	return string(data), nil
}
