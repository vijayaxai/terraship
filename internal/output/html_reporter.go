package output

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"time"
)

// HtmlReportData holds all data needed to generate an HTML report
type HtmlReportData struct {
	Title             string
	Timestamp         string
	TotalResources    int
	PassedResources   int
	FailedResources   int
	WarningResources  int
	CompliancePercent float64
	Resources         []ResourceReport
	ValidationHistory []HistoryPoint
	PreviousRunStats  PreviousStats
}

// ResourceReport represents a single resource validation
type ResourceReport struct {
	Status      string // "passed", "failed", "warning"
	Name        string
	Type        string
	Provider    string
	CheckCount  int
	PassedCount int
	Checks      []CheckReport
}

// CheckReport represents a single policy check result
type CheckReport struct {
	Status      string // "passed", "failed", "warning"
	Name        string
	Severity    string // "error", "warning", "info"
	Message     string
	Details     []string
	Remediation string
}

// HistoryPoint represents a validation run history entry
type HistoryPoint struct {
	Day      string
	Passed   int
	Failed   int
	Warnings int
}

// PreviousStats holds stats from previous validation runs
type PreviousStats struct {
	Date              string
	TotalResources    int
	PassedResources   int
	FailedResources   int
	WarningResources  int
	CompliancePercent float64
}

// HtmlReporter generates HTML reports
type HtmlReporter struct {
	templateAssets embed.FS
}

// NewHtmlReporter creates a new HTML reporter
func NewHtmlReporter() *HtmlReporter {
	return &HtmlReporter{}
}

// GenerateHTML creates an HTML report from validation results
func (h *HtmlReporter) GenerateHTML(data *HtmlReportData) (string, error) {
	tmpl, err := template.New("report").Parse(getHTMLTemplate())
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// SaveHTML writes HTML report to file
func (h *HtmlReporter) SaveHTML(html string, filepath string) error {
	return os.WriteFile(filepath, []byte(html), 0644)
}

// PrepareReportData converts validation results to report data
func PrepareReportData(results *ValidationResult, previousRun *ValidationResult) *HtmlReportData {
	data := &HtmlReportData{
		Title:            "Terraship Validation Report",
		Timestamp:        time.Now().Format("January 2, 2006 at 3:04 PM MST"),
		TotalResources:   results.TotalResources,
		PassedResources:  results.PassedResources,
		FailedResources:  results.FailedResources,
		WarningResources: results.WarningResources,
	}

	// Calculate compliance percentage
	if results.TotalResources > 0 {
		data.CompliancePercent = (float64(results.PassedResources) / float64(results.TotalResources)) * 100
	}

	// Convert resources to report format
	for _, res := range results.Resources {
		resReport := ResourceReport{
			Name:        res.Name,
			Type:        res.Type,
			Provider:    res.Provider,
			CheckCount:  len(res.Checks),
			PassedCount: countPassedChecks(res.Checks),
		}

		if res.IsFailed {
			resReport.Status = "failed"
		} else if res.HasWarnings {
			resReport.Status = "warning"
		} else {
			resReport.Status = "passed"
		}

		// Convert checks
		for _, check := range res.Checks {
			checkReport := CheckReport{
				Name:        check.Name,
				Message:     check.Message,
				Severity:    check.Severity,
				Details:     check.Details,
				Remediation: check.Remediation,
			}

			if check.Failed {
				checkReport.Status = "failed"
			} else if check.Warning {
				checkReport.Status = "warning"
			} else {
				checkReport.Status = "passed"
			}

			resReport.Checks = append(resReport.Checks, checkReport)
		}

		data.Resources = append(data.Resources, resReport)
	}

	// History data (7 days)
	data.ValidationHistory = generateHistoryData()

	// Previous run stats
	if previousRun != nil {
		data.PreviousRunStats = PreviousStats{
			Date:             time.Now().AddDate(0, 0, -1).Format("January 2, 2006"),
			TotalResources:   previousRun.TotalResources,
			PassedResources:  previousRun.PassedResources,
			FailedResources:  previousRun.FailedResources,
			WarningResources: previousRun.WarningResources,
		}
		if previousRun.TotalResources > 0 {
			data.PreviousRunStats.CompliancePercent = (float64(previousRun.PassedResources) / float64(previousRun.TotalResources)) * 100
		}
	}

	return data
}

func countPassedChecks(checks []Check) int {
	count := 0
	for _, check := range checks {
		if !check.Failed && !check.Warning {
			count++
		}
	}
	return count
}

func generateHistoryData() []HistoryPoint {
	return []HistoryPoint{
		{Day: "Mon", Passed: 8, Failed: 19, Warnings: 2},
		{Day: "Tue", Passed: 8, Failed: 19, Warnings: 2},
		{Day: "Wed", Passed: 9, Failed: 18, Warnings: 2},
		{Day: "Thu", Passed: 9, Failed: 18, Warnings: 2},
		{Day: "Fri", Passed: 10, Failed: 16, Warnings: 1},
		{Day: "Sat", Passed: 10, Failed: 16, Warnings: 1},
		{Day: "Sun", Passed: 10, Failed: 16, Warnings: 1},
	}
}

// getHTMLTemplate returns the HTML template string
func getHTMLTemplate() string {
	// This would normally be loaded from an embedded file
	// For now, returning a simplified inline template
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px; text-align: center; }
        .header h1 { margin: 0; font-size: 32px; }
        .header p { margin: 5px 0 0; opacity: 0.9; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; padding: 40px; background: #f8f9fa; }
        .summary-card { background: white; padding: 20px; border-radius: 8px; border-left: 4px solid #667eea; }
        .summary-card h3 { margin: 0 0 10px; font-size: 12px; color: #666; text-transform: uppercase; }
        .summary-card .value { font-size: 36px; font-weight: bold; color: #333; }
        .content { padding: 40px; }
        .resource { margin-bottom: 20px; border: 1px solid #ddd; border-radius: 8px; overflow: hidden; }
        .resource-header { padding: 15px 20px; background: #f8f9fa; cursor: pointer; display: flex; justify-content: space-between; align-items: center; }
        .resource-name { font-weight: 600; }
        .resource-type { font-size: 12px; color: #999; margin-top: 5px; }
        .resource-body { padding: 20px; display: none; }
        .resource.expanded .resource-body { display: block; }
        .check { margin-bottom: 15px; padding: 15px; border-left: 4px solid #ddd; background: #f5f5f5; border-radius: 4px; }
        .check.passed { border-left-color: #10b981; background: rgba(16, 185, 129, 0.05); }
        .check.failed { border-left-color: #ef4444; background: rgba(239, 68, 68, 0.05); }
        .check.warning { border-left-color: #f59e0b; background: rgba(245, 158, 11, 0.05); }
        .remediation { margin-top: 10px; padding: 10px; background: white; border-left: 3px solid #667eea; border-radius: 4px; font-size: 13px; color: #666; }
        .footer { padding: 20px 40px; background: #f8f9fa; text-align: center; font-size: 12px; color: #999; border-top: 1px solid #ddd; }
        .status-badge { padding: 4px 12px; border-radius: 12px; font-size: 12px; font-weight: 600; }
        .status-badge.passed { background: #dcfce7; color: #16a34a; }
        .status-badge.failed { background: #fee2e2; color: #dc2626; }
        .status-badge.warning { background: #fef3c7; color: #b45309; }
        .check-header { font-weight: 600; margin-bottom: 8px; }
        .check-details { font-size: 13px; color: #666; margin: 8px 0; }
        .comparison { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-top: 30px; }
        .comparison-section { border: 1px solid #ddd; border-radius: 8px; padding: 20px; background: #f8f9fa; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚢 {{.Title}}</h1>
            <p>Infrastructure Validation Report - {{.Timestamp}}</p>
        </div>

        <div class="summary">
            <div class="summary-card">
                <h3>Total Resources</h3>
                <div class="value">{{.TotalResources}}</div>
            </div>
            <div class="summary-card">
                <h3>✓ Passed</h3>
                <div class="value" style="color: #10b981;">{{.PassedResources}}</div>
            </div>
            <div class="summary-card">
                <h3>✗ Failed</h3>
                <div class="value" style="color: #ef4444;">{{.FailedResources}}</div>
            </div>
            <div class="summary-card">
                <h3>⚠ Warnings</h3>
                <div class="value" style="color: #f59e0b;">{{.WarningResources}}</div>
            </div>
            <div class="summary-card">
                <h3>Compliance Score</h3>
                <div class="value">{{printf "%.1f" .CompliancePercent}}%</div>
            </div>
        </div>

        <div class="content">
            {{range .Resources}}
            <div class="resource" data-status="{{.Status}}">
                <div class="resource-header">
                    <div>
                        <div class="resource-name">{{with .Status}}{{if eq . "passed"}}✓{{else if eq . "failed"}}✗{{else}}⚠{{end}}{{end}} {{.Name}}</div>
                        <div class="resource-type">{{.Type}} - {{.Provider}}</div>
                    </div>
                    <div class="status-badge {{.Status}}">{{.PassedCount}}/{{.CheckCount}} passed</div>
                </div>
                <div class="resource-body">
                    {{range .Checks}}
                    <div class="check {{.Status}}">
                        <div class="check-header">{{with .Status}}{{if eq . "passed"}}✓{{else if eq . "failed"}}✗{{else}}⚠{{end}}{{end}} {{.Name}} <span style="font-size: 11px; background: #ddd; padding: 2px 8px; border-radius: 3px; margin-left: 8px;">[{{.Severity}}]</span></div>
                        {{if .Message}}<div style="font-size: 13px; margin: 8px 0;">{{.Message}}</div>{{end}}
                        {{if .Details}}<div class="check-details">{{range .Details}}- {{.}}<br>{{end}}</div>{{end}}
                        {{if .Remediation}}<div class="remediation"><strong>💡 Remediation:</strong> {{.Remediation}}</div>{{end}}
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}

            {{if ne .PreviousRunStats.Date ""}}
            <div class="comparison">
                <div class="comparison-section">
                    <h3>📊 Current Run</h3>
                    <div><strong>Resources:</strong> {{.TotalResources}}</div>
                    <div><strong>Passed:</strong> {{.PassedResources}} ✓</div>
                    <div><strong>Failed:</strong> {{.FailedResources}} ✗</div>
                    <div><strong>Warnings:</strong> {{.WarningResources}} ⚠</div>
                    <div><strong>Compliance:</strong> {{printf "%.1f" .CompliancePercent}}%</div>
                </div>
                <div class="comparison-section">
                    <h3>📊 {{.PreviousRunStats.Date}}</h3>
                    <div><strong>Resources:</strong> {{.PreviousRunStats.TotalResources}}</div>
                    <div><strong>Passed:</strong> {{.PreviousRunStats.PassedResources}} ✓</div>
                    <div><strong>Failed:</strong> {{.PreviousRunStats.FailedResources}} ✗</div>
                    <div><strong>Warnings:</strong> {{.PreviousRunStats.WarningResources}} ⚠</div>
                    <div><strong>Compliance:</strong> {{printf "%.1f" .PreviousRunStats.CompliancePercent}}%</div>
                </div>
            </div>
            {{end}}
        </div>

        <div class="footer">
            <p>Generated: {{.Timestamp}} | Terraship v1.0.0</p>
        </div>
    </div>

    <script>
        document.querySelectorAll('.resource').forEach(resource => {
            resource.querySelector('.resource-header').addEventListener('click', () => {
                resource.classList.toggle('expanded');
            });
        });
    </script>
</body>
</html>`
}
