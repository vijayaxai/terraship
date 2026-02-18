// Package cloud provides the interface for multi-cloud adapters.
package cloud

import (
	"context"
)

// Provider represents supported cloud providers
type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderAzure Provider = "azure"
	ProviderGCP   Provider = "gcp"
	ProviderNone  Provider = "none"
)

// ResourceStatus represents the status of a cloud resource
type ResourceStatus struct {
	ResourceID   string                 `json:"resource_id"`
	ResourceType string                 `json:"resource_type"`
	Exists       bool                   `json:"exists"`
	State        string                 `json:"state,omitempty"`
	Tags         map[string]string      `json:"tags,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	DriftDetected bool                  `json:"drift_detected"`
	DriftDetails  []string              `json:"drift_details,omitempty"`
}

// ValidationResult contains the result of a resource validation
type ValidationResult struct {
	ResourceID  string   `json:"resource_id"`
	RuleName    string   `json:"rule_name"`
	Passed      bool     `json:"passed"`
	Message     string   `json:"message"`
	Severity    string   `json:"severity"` // "error", "warning", "info"
	Remediation string   `json:"remediation,omitempty"`
	Details     []string `json:"details,omitempty"`
}

// CloudConfig contains configuration for cloud provider authentication
type CloudConfig struct {
	Provider Provider
	Region   string

	// AWS specific
	AWSProfile  string
	AWSRoleARN  string
	AWSRegion   string

	// Azure specific
	AzureSubscriptionID string
	AzureTenantID       string
	AzureClientID       string
	AzureClientSecret   string

	// GCP specific
	GCPProject            string
	GCPCredentialsFile    string
	GCPServiceAccountJSON string
}

// Adapter defines the interface for cloud provider operations
type Adapter interface {
	// Name returns the cloud provider name
	Name() Provider

	// Initialize sets up the cloud adapter with credentials
	Initialize(ctx context.Context, config CloudConfig) error

	// DetectProvider attempts to auto-detect if this is the correct provider
	// based on environment variables, credentials, or resource naming
	DetectProvider(ctx context.Context) (bool, float64, error) // bool=detected, float64=confidence (0.0-1.0)

	// ValidateCredentials checks if the configured credentials are valid
	ValidateCredentials(ctx context.Context) error

	// GetResourceStatus retrieves the current status of a resource from the cloud
	GetResourceStatus(ctx context.Context, resourceType, resourceID string) (*ResourceStatus, error)

	// ValidateResourceCompliance checks if a resource complies with policies
	ValidateResourceCompliance(ctx context.Context, resourceType string, resource map[string]interface{}, rules []ValidationRule) ([]ValidationResult, error)

	// DetectDrift compares Terraform state with actual cloud resources
	DetectDrift(ctx context.Context, plannedState map[string]interface{}, resourceType, resourceID string) (*ResourceStatus, error)

	// ListResources lists resources of a given type (optional; for discovery)
	ListResources(ctx context.Context, resourceType string) ([]string, error)

	// Close cleans up any resources used by the adapter
	Close() error
}

// ValidationRule represents a policy rule to validate
type ValidationRule struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description"`
	Severity    string                 `yaml:"severity" json:"severity"` // "error", "warning", "info"
	Category    string                 `yaml:"category" json:"category"` // "security", "compliance", "cost", "performance"
	Enabled     bool                   `yaml:"enabled" json:"enabled"`
	ResourceTypes []string             `yaml:"resource_types" json:"resource_types"`
	Conditions  map[string]interface{} `yaml:"conditions" json:"conditions"`
	Message     string                 `yaml:"message" json:"message"`
	Remediation string                 `yaml:"remediation" json:"remediation"`
}

// DetectionResult holds the result of provider auto-detection
type DetectionResult struct {
	Provider   Provider
	Confidence float64
	Reason     string
}

// AutoDetect attempts to detect the cloud provider from environment and context
func AutoDetect(ctx context.Context, adapters []Adapter) (*DetectionResult, error) {
	var bestMatch *DetectionResult

	for _, adapter := range adapters {
		detected, confidence, err := adapter.DetectProvider(ctx)
		if err != nil {
			continue
		}

		if detected && (bestMatch == nil || confidence > bestMatch.Confidence) {
			bestMatch = &DetectionResult{
				Provider:   adapter.Name(),
				Confidence: confidence,
				Reason:     "Auto-detected from environment",
			}
		}
	}

	if bestMatch == nil {
		return &DetectionResult{
			Provider:   ProviderNone,
			Confidence: 0.0,
			Reason:     "No cloud provider detected",
		}, nil
	}

	return bestMatch, nil
}
