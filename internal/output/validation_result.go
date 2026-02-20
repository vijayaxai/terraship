package output

import (
	"encoding/json"
	"time"
)

// ValidationResult holds complete validation results
type ValidationResult struct {
	TotalResources   int
	PassedResources  int
	FailedResources  int
	WarningResources int
	Timestamp        string
	Resources        []Resource
}

// Resource represents a validated resource
type Resource struct {
	Name        string
	Type        string
	Provider    string
	IsFailed    bool
	HasWarnings bool
	Checks      []Check
}

// Check represents a single policy check
type Check struct {
	Name        string
	Message     string
	Severity    string // "error", "warning", "info"
	Failed      bool
	Warning     bool
	Details     []string
	Remediation string
}

// ToJSON converts results to JSON
func (vr *ValidationResult) ToJSON() ([]byte, error) {
	data := map[string]interface{}{
		"timestamp":          vr.Timestamp,
		"total_resources":    vr.TotalResources,
		"passed_resources":   vr.PassedResources,
		"failed_resources":   vr.FailedResources,
		"warning_resources":  vr.WarningResources,
		"compliance_percent": calculateCompliance(vr.TotalResources, vr.PassedResources),
		"resources":          vr.Resources,
		"validation_passed":  vr.FailedResources == 0,
	}

	return json.MarshalIndent(data, "", "  ")
}

// ToSARIF converts results to SARIF format
func (vr *ValidationResult) ToSARIF() ([]byte, error) {
	// SARIF 2.1.0 format for GitHub Code Scanning and other tools
	sarifResults := map[string]interface{}{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":           "Terraship",
						"version":        "1.0.0",
						"informationUri": "https://github.com/vijayaxai/terraship",
					},
				},
				"results": buildSARIFResults(vr),
			},
		},
	}

	return json.MarshalIndent(sarifResults, "", "  ")
}

// buildSARIFResults converts validation results to SARIF format
func buildSARIFResults(vr *ValidationResult) []map[string]interface{} {
	var results []map[string]interface{}

	for _, resource := range vr.Resources {
		for _, check := range resource.Checks {
			if check.Failed || check.Warning {
				level := "warning"
				if check.Failed {
					level = "error"
				}

				result := map[string]interface{}{
					"ruleId": check.Name,
					"level":  level,
					"message": map[string]interface{}{
						"text": check.Message,
					},
					"locations": []map[string]interface{}{
						{
							"physicalLocation": map[string]interface{}{
								"artifactLocation": map[string]interface{}{
									"uri": resource.Name,
								},
							},
						},
					},
					"properties": map[string]interface{}{
						"resource_type": resource.Type,
						"provider":      resource.Provider,
						"severity":      check.Severity,
					},
				}

				results = append(results, result)
			}
		}
	}

	return results
}

// CalculateCompliance returns compliance percentage
func calculateCompliance(total, passed int) float64 {
	if total == 0 {
		return 0
	}
	return (float64(passed) / float64(total)) * 100
}

// ComparisonReport represents a comparison between two validation runs
type ComparisonReport struct {
	Current          *ValidationResult
	Previous         *ValidationResult
	ChangedResources []ResourceChange
	TrendPercent     float64 // positive = improving, negative = regressing
}

// ResourceChange tracks changes in a resource validation
type ResourceChange struct {
	ResourceName   string
	Status         string // "improved", "regressed", "unchanged"
	PreviousFailed int
	CurrentFailed  int
	DetailChanges  []string
}

// Compare compares two validation runs
func Compare(current, previous *ValidationResult) *ComparisonReport {
	report := &ComparisonReport{
		Current:  current,
		Previous: previous,
	}

	// Calculate trend
	if previous != nil {
		prevCompliance := calculateCompliance(previous.TotalResources, previous.PassedResources)
		currCompliance := calculateCompliance(current.TotalResources, current.PassedResources)
		report.TrendPercent = currCompliance - prevCompliance
	}

	return report
}

// ExportStats returns exportable statistics
type ExportStats struct {
	Timestamp           time.Time      `json:"timestamp"`
	TotalResources      int            `json:"total_resources"`
	PassedResources     int            `json:"passed_resources"`
	FailedResources     int            `json:"failed_resources"`
	WarningResources    int            `json:"warning_resources"`
	CompliancePercent   float64        `json:"compliance_percent"`
	ResourcesByType     map[string]int `json:"resources_by_type"`
	ResourcesByProvider map[string]int `json:"resources_by_provider"`
	FailuresByRule      map[string]int `json:"failures_by_rule"`
}

// GetExportStats generates exportable statistics from results
func (vr *ValidationResult) GetExportStats() *ExportStats {
	stats := &ExportStats{
		Timestamp:           time.Now(),
		TotalResources:      vr.TotalResources,
		PassedResources:     vr.PassedResources,
		FailedResources:     vr.FailedResources,
		WarningResources:    vr.WarningResources,
		CompliancePercent:   calculateCompliance(vr.TotalResources, vr.PassedResources),
		ResourcesByType:     make(map[string]int),
		ResourcesByProvider: make(map[string]int),
		FailuresByRule:      make(map[string]int),
	}

	for _, resource := range vr.Resources {
		stats.ResourcesByType[resource.Type]++
		stats.ResourcesByProvider[resource.Provider]++

		for _, check := range resource.Checks {
			if check.Failed {
				stats.FailuresByRule[check.Name]++
			}
		}
	}

	return stats
}
