// Package core orchestrates the validation workflow.
package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vijayaxai/terraship/internal/cloud"
	awsadapter "github.com/vijayaxai/terraship/internal/cloud/aws"
	azureadapter "github.com/vijayaxai/terraship/internal/cloud/azure"
	gcpadapter "github.com/vijayaxai/terraship/internal/cloud/gcp"
	"github.com/vijayaxai/terraship/internal/rules"
	"github.com/vijayaxai/terraship/internal/terraform"
)

// ValidationMode defines how validation is performed
type ValidationMode string

const (
	// ModeValidateExisting validates existing infrastructure without apply
	ModeValidateExisting ValidationMode = "validate-existing"
	// ModeEphemeralSandbox creates temporary infrastructure for testing
	ModeEphemeralSandbox ValidationMode = "ephemeral-sandbox"
)

// ValidatorConfig holds configuration for the validator
type ValidatorConfig struct {
	Mode          ValidationMode
	WorkingDir    string
	PolicyPath    string
	CloudProvider string // manual override; empty for auto-detect
	OutputFormat  string // "human", "json", "sarif"
	OutputFile    string
	NoDestroy     bool // for ephemeral mode
	Verbose       bool
}

// Validator orchestrates the validation process
type Validator struct {
	config       ValidatorConfig
	tfClient     *terraform.Client
	cloudAdapter cloud.Adapter
	rulesEngine  *rules.Engine
	results      []ValidationReport
}

// ValidationReport contains the results of validation
type ValidationReport struct {
	ResourceAddress string                   `json:"resource_address"`
	ResourceType    string                   `json:"resource_type"`
	Provider        string                   `json:"provider"`
	Status          string                   `json:"status"` // "pass", "fail", "warning", "error"
	RuleResults     []cloud.ValidationResult `json:"rule_results"`
	DriftStatus     *cloud.ResourceStatus    `json:"drift_status,omitempty"`
	Errors          []string                 `json:"errors,omitempty"`
}

// Summary provides overall validation summary
type Summary struct {
	TotalResources   int                `json:"total_resources"`
	PassedResources  int                `json:"passed_resources"`
	FailedResources  int                `json:"failed_resources"`
	WarningResources int                `json:"warning_resources"`
	ErrorResources   int                `json:"error_resources"`
	DriftDetected    int                `json:"drift_detected"`
	Reports          []ValidationReport `json:"reports"`
}

// NewValidator creates a new validator instance
func NewValidator(config ValidatorConfig) (*Validator, error) {
	// Validate config
	if config.WorkingDir == "" {
		return nil, fmt.Errorf("working directory is required")
	}

	if config.PolicyPath == "" {
		return nil, fmt.Errorf("policy path is required")
	}

	// Create Terraform client
	tfClient, err := terraform.NewClient(config.WorkingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create terraform client: %w", err)
	}

	// Load rules engine
	rulesEngine, err := rules.NewEngine(config.PolicyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &Validator{
		config:      config,
		tfClient:    tfClient,
		rulesEngine: rulesEngine,
		results:     make([]ValidationReport, 0),
	}, nil
}

// Validate performs the validation workflow
func (v *Validator) Validate(ctx context.Context) (*Summary, error) {
	// Step 1: Initialize Terraform
	if err := v.tfClient.Init(ctx, false); err != nil {
		return nil, fmt.Errorf("terraform init failed: %w", err)
	}

	// Step 2: Validate Terraform configuration
	if err := v.tfClient.Validate(ctx); err != nil {
		return nil, fmt.Errorf("terraform validate failed: %w", err)
	}

	// Step 3: Detect or set cloud provider
	provider := v.config.CloudProvider
	if provider == "" {
		detectedProvider, err := v.tfClient.GetProvider(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to detect cloud provider: %w", err)
		}
		provider = detectedProvider
	}

	// Step 4: Initialize cloud adapter
	if err := v.initializeCloudAdapter(ctx, provider); err != nil {
		return nil, fmt.Errorf("failed to initialize cloud adapter: %w", err)
	}

	// Step 5: Generate Terraform plan
	planFile := filepath.Join(os.TempDir(), "terraship-plan.tfplan")
	defer os.Remove(planFile)

	if err := v.tfClient.Plan(ctx, planFile); err != nil {
		return nil, fmt.Errorf("terraform plan failed: %w", err)
	}

	// Step 6: Parse plan output
	plan, err := v.tfClient.ShowJSON(ctx, planFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	// Step 7: Validate resources
	if err := v.validateResources(ctx, plan); err != nil {
		return nil, fmt.Errorf("resource validation failed: %w", err)
	}

	// Step 8: For ephemeral mode, apply and then destroy
	if v.config.Mode == ModeEphemeralSandbox {
		if err := v.runEphemeralMode(ctx, planFile); err != nil {
			return nil, fmt.Errorf("ephemeral mode failed: %w", err)
		}
	}

	// Step 9: Generate summary
	summary := v.generateSummary()

	return summary, nil
}

func (v *Validator) initializeCloudAdapter(ctx context.Context, provider string) error {
	var adapter cloud.Adapter

	switch provider {
	case "aws":
		adapter = newAWSAdapter() // Will implement next
	case "azure":
		adapter = newAzureAdapter()
	case "gcp":
		adapter = newGCPAdapter()
	default:
		return fmt.Errorf("unsupported cloud provider: %s", provider)
	}

	config := cloud.CloudConfig{
		Provider: cloud.Provider(provider),
	}

	if err := adapter.Initialize(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize %s adapter: %w", provider, err)
	}

	if err := adapter.ValidateCredentials(ctx); err != nil {
		return fmt.Errorf("cloud credentials validation failed: %w", err)
	}

	v.cloudAdapter = adapter
	return nil
}

func (v *Validator) validateResources(ctx context.Context, plan *terraform.PlanOutput) error {
	if plan.PlannedValues == nil || plan.PlannedValues.RootModule == nil {
		return fmt.Errorf("no resources found in plan")
	}

	// Collect all resources from root and child modules
	resources := v.collectResources(plan.PlannedValues.RootModule)

	for _, resource := range resources {
		report := v.validateResource(ctx, resource)
		v.results = append(v.results, report)
	}

	return nil
}

func (v *Validator) collectResources(module *terraform.Module) []terraform.Resource {
	var resources []terraform.Resource

	resources = append(resources, module.Resources...)

	for _, childModule := range module.ChildModules {
		resources = append(resources, v.collectResources(&childModule)...)
	}

	return resources
}

func (v *Validator) validateResource(ctx context.Context, resource terraform.Resource) ValidationReport {
	report := ValidationReport{
		ResourceAddress: resource.Address,
		ResourceType:    resource.Type,
		Provider:        resource.ProviderName,
		Status:          "pass",
		RuleResults:     make([]cloud.ValidationResult, 0),
		Errors:          make([]string, 0),
	}

	// Get applicable rules
	applicableRules := v.rulesEngine.GetRulesForResource(resource.Type)

	// Evaluate each rule
	for _, rule := range applicableRules {
		result := v.rulesEngine.EvaluateRule(rule, resource.Values)
		result.ResourceID = resource.Address
		report.RuleResults = append(report.RuleResults, result)

		if !result.Passed {
			if result.Severity == "error" {
				report.Status = "fail"
			} else if result.Severity == "warning" && report.Status != "fail" {
				report.Status = "warning"
			}
		}
	}

	// Check for drift if in validate-existing mode
	if v.config.Mode == ModeValidateExisting && v.cloudAdapter != nil {
		resourceID := v.extractResourceID(resource)
		if resourceID != "" {
			driftStatus, err := v.cloudAdapter.DetectDrift(ctx, resource.Values, resource.Type, resourceID)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Drift detection failed: %s", err))
			} else {
				report.DriftStatus = driftStatus
				if driftStatus.DriftDetected {
					if report.Status == "pass" {
						report.Status = "warning"
					}
				}
			}
		}
	}

	return report
}

func (v *Validator) extractResourceID(resource terraform.Resource) string {
	// Try to extract resource ID from values
	if id, ok := resource.Values["id"].(string); ok && id != "" {
		return id
	}
	if name, ok := resource.Values["name"].(string); ok && name != "" {
		return name
	}
	if arn, ok := resource.Values["arn"].(string); ok && arn != "" {
		return arn
	}
	return ""
}

func (v *Validator) runEphemeralMode(ctx context.Context, planFile string) error {
	// Apply the plan
	applyErr := v.tfClient.Apply(ctx, planFile)

	// Always attempt destroy unless --no-destroy flag is set, even if apply failed
	// This ensures cleanup happens to prevent resource leaks
	if !v.config.NoDestroy {
		if err := v.tfClient.Destroy(ctx, true); err != nil {
			// Log warning but don't block error reporting from apply failure
			if v.config.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: terraform destroy encountered issues: %v\n", err)
			}
		}
	}

	// Return the apply error after cleanup attempt
	if applyErr != nil {
		return fmt.Errorf("terraform apply failed: %w", applyErr)
	}

	return nil
}

func (v *Validator) generateSummary() *Summary {
	summary := &Summary{
		TotalResources: len(v.results),
		Reports:        v.results,
	}

	for _, report := range v.results {
		switch report.Status {
		case "pass":
			summary.PassedResources++
		case "fail":
			summary.FailedResources++
		case "warning":
			summary.WarningResources++
		case "error":
			summary.ErrorResources++
		}

		if report.DriftStatus != nil && report.DriftStatus.DriftDetected {
			summary.DriftDetected++
		}
	}

	return summary
}

// Cloud adapter factory functions
func newAWSAdapter() cloud.Adapter {
	return awsadapter.NewAdapter()
}

func newAzureAdapter() cloud.Adapter {
	return azureadapter.NewAdapter()
}

func newGCPAdapter() cloud.Adapter {
	return gcpadapter.NewAdapter()
}
