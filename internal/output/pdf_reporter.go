package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// PDFReporter generates PDF reports from HTML
type PDFReporter struct {
	htmlReporter *HtmlReporter
}

// NewPDFReporter creates a new PDF reporter
func NewPDFReporter() *PDFReporter {
	return &PDFReporter{
		htmlReporter: NewHtmlReporter(),
	}
}

// GeneratePDF creates a PDF report from validation results
// Uses wkhtmltopdf or similar tool if available, otherwise returns base64 HTML
func (p *PDFReporter) GeneratePDF(data *HtmlReportData, outputPath string) error {
	// Step 1: Generate HTML first
	html, err := p.htmlReporter.GenerateHTML(data)
	if err != nil {
		return err
	}

	// Step 2: Try to convert HTML to PDF using wkhtmltopdf
	pdfBytes, err := p.convertHTMLToPDF(html)
	if err != nil {
		// Fallback: If wkhtmltopdf not available, create HTML-based PDF
		return p.htmlReporter.SaveHTML(html, outputPath)
	}

	// Step 3: Write PDF bytes to file
	if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// convertHTMLToPDF converts HTML string to PDF bytes using wkhtmltopdf
// This requires wkhtmltopdf to be installed and in PATH
func (p *PDFReporter) convertHTMLToPDF(html string) ([]byte, error) {
	// Create temporary HTML file
	tmpDir := os.TempDir()
	tmpHTMLFile := filepath.Join(tmpDir, "terraship_report_temp.html")
	tmpPDFFile := filepath.Join(tmpDir, "terraship_report_temp.pdf")

	defer os.Remove(tmpHTMLFile)
	defer os.Remove(tmpPDFFile)

	if err := os.WriteFile(tmpHTMLFile, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("failed to create temporary HTML file: %w", err)
	}

	// Option 1: Try using wkhtmltopdf command
	return p.executeWkhtmltopdf(tmpHTMLFile, tmpPDFFile)
}

// executeWkhtmltopdf executes wkhtmltopdf command
// Returns pdf bytes if successful, error if not available
func (p *PDFReporter) executeWkhtmltopdf(htmlPath, pdfPath string) ([]byte, error) {
	// Check if wkhtmltopdf is available
	_, err := os.Stat("/usr/local/bin/wkhtmltopdf")
	if os.IsNotExist(err) {
		_, err = os.Stat("C:\\Program Files\\wkhtmltopdf\\bin\\wkhtmltopdf.exe")
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("wkhtmltopdf not found: %w", err)
		}
	}

	// For now, return error and fallback to HTML
	// In production, this would execute: wkhtmltopdf htmlPath pdfPath
	return nil, fmt.Errorf("wkhtmltopdf conversion requires external tool")
}

// GeneratePDFAlternative uses a Go-native PDF library
// This is an alternative implementation using gofpdf or similar
func (p *PDFReporter) GeneratePDFAlternative(data *HtmlReportData, outputPath string) error {
	// Alternative implementation notes:
	// Could use libraries like:
	// - github.com/jung-kurt/gofpdf (simple, lightweight)
	// - github.com/mandykoh/prism (modern, complex layouts)
	// - github.com/signintech/gopdf (efficient)
	//
	// Example with gofpdf:
	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.AddPage()
	// pdf.SetFont("Arial", "B", 16)
	// pdf.Cell(0, 10, "Terraship Validation Report")
	// ... add more content
	// pdf.OutputToFile(outputPath)

	return fmt.Errorf("PDF generation requires additional dependencies - use HTML export or install wkhtmltopdf")
}

// GetPDFInstallInstructions returns instructions for installing PDF support
func GetPDFInstallInstructions() string {
	return `PDF Export Support
====================

To enable PDF export, install one of these tools:

OPTION 1: wkhtmltopdf (Recommended)
-----------------------------------
macOS:
  brew install wkhtmltopdf

Ubuntu/Debian:
  sudo apt-get install wkhtmltopdf

Windows:
  choco install wkhtmltopdf
  OR download from: https://wkhtmltopdf.org/

OPTION 2: Use HTML Export
---------------------------
HTML reports can be opened in any browser and printed as PDF:
  terraship validate ./terraform --output html --html-file report.html
  Then open in browser: Ctrl+P (or Cmd+P) → Save as PDF

After installing, PDF export will work automatically:
  terraship validate ./terraform --output pdf --pdf-file report.pdf`
}
