package output

import (
	"encoding/json"
	"testing"
)

// TestValidationResult_ToJSON tests JSON export
func TestValidationResult_ToJSON(t *testing.T) {
	result := &ValidationResult{
		TotalResources:   10,
		PassedResources:  8,
		FailedResources:  2,
		WarningResources: 0,
		Timestamp:        "2026-02-19 11:15 AM",
		Resources: []Resource{
			{
				Name:     "aws_s3_bucket_example",
				Type:     "aws_s3_bucket",
				Provider: "aws",
			},
		},
	}

	jsonBytes, err := result.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}

	// Verify JSON is valid
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	// Verify expected fields
	if data["total_resources"] != float64(10) {
		t.Errorf("Expected total_resources=10, got %v", data["total_resources"])
	}
	if data["passed_resources"] != float64(8) {
		t.Errorf("Expected passed_resources=8, got %v", data["passed_resources"])
	}
	if data["failed_resources"] != float64(2) {
		t.Errorf("Expected failed_resources=2, got %v", data["failed_resources"])
	}
}

// TestValidationResult_ToSARIF tests SARIF export format
func TestValidationResult_ToSARIF(t *testing.T) {
	result := &ValidationResult{
		TotalResources:   10,
		PassedResources:  8,
		FailedResources:  2,
		WarningResources: 0,
		Timestamp:        "2026-02-19 11:15 AM",
		Resources: []Resource{
			{
				Name:     "aws_s3_bucket_example",
				Type:     "aws_s3_bucket",
				Provider: "aws",
				Checks: []Check{
					{
						Name:     "encryption_at_rest",
						Severity: "error",
						Message:  "Encryption not enabled",
						Failed:   true,
					},
				},
			},
		},
	}

	sarifBytes, err := result.ToSARIF()
	if err != nil {
		t.Fatalf("ToSARIF() failed: %v", err)
	}

	// Verify SARIF is valid JSON
	var sarif map[string]interface{}
	if err := json.Unmarshal(sarifBytes, &sarif); err != nil {
		t.Fatalf("Invalid SARIF output: %v", err)
	}

	// Verify SARIF version
	if sarif["version"] != "2.1.0" {
		t.Errorf("Expected SARIF version 2.1.0, got %v", sarif["version"])
	}

	// Verify runs exist
	if runs, ok := sarif["runs"]; !ok || runs == nil {
		t.Error("SARIF output missing 'runs' field")
	}
}

// TestValidationResult_Compliance calculates compliance percentage
func TestValidationResult_Compliance(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		passed   int
		expected float64
	}{
		{"Perfect compliance", 10, 10, 100.0},
		{"80% compliance", 10, 8, 80.0},
		{"50% compliance", 10, 5, 50.0},
		{"Zero resources", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ValidationResult{
				TotalResources:  tt.total,
				PassedResources: tt.passed,
			}

			compliance := 0.0
			if result.TotalResources > 0 {
				compliance = (float64(result.PassedResources) / float64(result.TotalResources)) * 100
			}

			if compliance != tt.expected {
				t.Errorf("Expected compliance %.1f%%, got %.1f%%", tt.expected, compliance)
			}
		})
	}
}

// TestHtmlReporter_GenerateHTML tests HTML report generation
func TestHtmlReporter_GenerateHTML(t *testing.T) {
	reporter := NewHtmlReporter()

	reportData := &HtmlReportData{
		Title:             "Test Report",
		Timestamp:         "2026-02-19 11:15 AM",
		TotalResources:    10,
		PassedResources:   8,
		FailedResources:   2,
		CompliancePercent: 80.0,
		Resources: []ResourceReport{
			{
				Name:        "test-resource",
				Type:        "aws_s3_bucket",
				Provider:    "aws",
				Status:      "passed",
				CheckCount:  5,
				PassedCount: 5,
			},
		},
	}

	html, err := reporter.GenerateHTML(reportData)
	if err != nil {
		t.Fatalf("GenerateHTML() failed: %v", err)
	}

	// Verify HTML contains expected elements
	if len(html) == 0 {
		t.Error("Generated HTML is empty")
	}

	expectedStrings := []string{
		"<!DOCTYPE html>", // Must be valid HTML
		"80",              // Compliance score
		"test-resource",
		"aws_s3_bucket",
	}

	for _, expected := range expectedStrings {
		if !contains(html, expected) {
			t.Errorf("HTML missing expected content: %s", expected)
		}
	}
}

// TestPDFReporter_Initialization tests PDF reporter initialization
func TestPDFReporter_Initialization(t *testing.T) {
	reporter := NewPDFReporter()
	if reporter == nil {
		t.Error("NewPDFReporter() returned nil")
	}

	// Verify it has the HTML reporter
	if reporter.htmlReporter == nil {
		t.Error("PDFReporter missing htmlReporter")
	}
}

// TestGetPDFInstallInstructions verifies installation help text
func TestGetPDFInstallInstructions(t *testing.T) {
	instructions := GetPDFInstallInstructions()

	if len(instructions) == 0 {
		t.Error("GetPDFInstallInstructions returned empty string")
	}

	// Should mention installation methods
	if !contains(instructions, "brew") && !contains(instructions, "apt") {
		t.Error("Instructions should include package manager details")
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
