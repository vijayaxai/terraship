// Package gcp implements the GCP cloud adapter.
package gcp

import (
	"context"
	"fmt"
	"os"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/storage"
	"github.com/vijayaxai/terraship/internal/cloud"
	"google.golang.org/api/option"
)

// Adapter implements cloud.Adapter for GCP
type Adapter struct {
	projectID       string
	computeClient   *compute.InstancesClient
	storageClient   *storage.Client
	credentialsFile string
}

// NewAdapter creates a new GCP adapter
func NewAdapter() *Adapter {
	return &Adapter{}
}

// Name returns the provider name
func (a *Adapter) Name() cloud.Provider {
	return cloud.ProviderGCP
}

// Initialize sets up the GCP adapter
func (a *Adapter) Initialize(ctx context.Context, cloudConfig cloud.CloudConfig) error {
	var err error
	var opts []option.ClientOption

	// Get project ID
	a.projectID = cloudConfig.GCPProject
	if a.projectID == "" {
		a.projectID = os.Getenv("GCP_PROJECT")
	}
	if a.projectID == "" {
		a.projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if a.projectID == "" {
		return fmt.Errorf("GCP project ID is required")
	}

	// Set credentials
	if cloudConfig.GCPCredentialsFile != "" {
		a.credentialsFile = cloudConfig.GCPCredentialsFile
		opts = append(opts, option.WithCredentialsFile(cloudConfig.GCPCredentialsFile))
	} else if credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credFile != "" {
		a.credentialsFile = credFile
		opts = append(opts, option.WithCredentialsFile(credFile))
	}

	// Initialize compute client
	a.computeClient, err = compute.NewInstancesRESTClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create GCP compute client: %w", err)
	}

	// Initialize storage client
	a.storageClient, err = storage.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create GCP storage client: %w", err)
	}

	return nil
}

// DetectProvider attempts to detect if GCP is the provider
func (a *Adapter) DetectProvider(ctx context.Context) (bool, float64, error) {
	confidence := 0.0

	// Check for GCP environment variables
	if os.Getenv("GCP_PROJECT") != "" || os.Getenv("GOOGLE_CLOUD_PROJECT") != "" {
		confidence += 0.3
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		confidence += 0.4
	}
	if os.Getenv("GCLOUD_PROJECT") != "" {
		confidence += 0.3
	}

	return confidence > 0.5, confidence, nil
}

// ValidateCredentials checks if GCP credentials are valid
func (a *Adapter) ValidateCredentials(ctx context.Context) error {
	// Try to list instances as a lightweight validation
	req := &computepb.AggregatedListInstancesRequest{
		Project:    a.projectID,
		MaxResults: func() *uint32 { v := uint32(1); return &v }(),
	}

	it := a.computeClient.AggregatedList(ctx, req)
	_, err := it.Next()
	if err != nil && err.Error() != "no more items in iterator" {
		return fmt.Errorf("GCP credentials validation failed: %w", err)
	}

	return nil
}

// GetResourceStatus retrieves the current status of a GCP resource
func (a *Adapter) GetResourceStatus(ctx context.Context, resourceType, resourceID string) (*cloud.ResourceStatus, error) {
	status := &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Properties:   make(map[string]interface{}),
	}

	switch {
	case strings.HasPrefix(resourceType, "google_compute_instance"):
		return a.getComputeInstanceStatus(ctx, resourceID)
	case strings.HasPrefix(resourceType, "google_storage_bucket"):
		return a.getStorageBucketStatus(ctx, resourceID)
	default:
		return status, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (a *Adapter) getComputeInstanceStatus(ctx context.Context, resourceID string) (*cloud.ResourceStatus, error) {
	// Parse resource ID: projects/{project}/zones/{zone}/instances/{name}
	parts := strings.Split(resourceID, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid GCP resource ID format")
	}

	zone := parts[3]
	instanceName := parts[5]

	req := &computepb.GetInstanceRequest{
		Project:  a.projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := a.computeClient.Get(ctx, req)
	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   resourceID,
			ResourceType: "google_compute_instance",
			Exists:       false,
		}, nil
	}

	tags := make(map[string]string)
	if instance.Labels != nil {
		for key, value := range instance.Labels {
			tags[key] = value
		}
	}

	status := &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: "google_compute_instance",
		Exists:       true,
		State:        instance.GetStatus(),
		Tags:         tags,
		Properties:   make(map[string]interface{}),
	}

	if instance.MachineType != nil {
		status.Properties["machine_type"] = *instance.MachineType
	}
	if instance.Zone != nil {
		status.Properties["zone"] = *instance.Zone
	}

	return status, nil
}

func (a *Adapter) getStorageBucketStatus(ctx context.Context, bucketName string) (*cloud.ResourceStatus, error) {
	bucket := a.storageClient.Bucket(bucketName)
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   bucketName,
			ResourceType: "google_storage_bucket",
			Exists:       false,
		}, nil
	}

	tags := make(map[string]string)
	if attrs.Labels != nil {
		for key, value := range attrs.Labels {
			tags[key] = value
		}
	}

	status := &cloud.ResourceStatus{
		ResourceID:   bucketName,
		ResourceType: "google_storage_bucket",
		Exists:       true,
		Tags:         tags,
		Properties:   make(map[string]interface{}),
	}

	status.Properties["location"] = attrs.Location
	status.Properties["storage_class"] = attrs.StorageClass
	status.Properties["versioning_enabled"] = attrs.VersioningEnabled

	// Check encryption
	if attrs.Encryption != nil {
		status.Properties["encryption_enabled"] = true
		status.Properties["encryption_key"] = attrs.Encryption.DefaultKMSKeyName
	}

	return status, nil
}

// ValidateResourceCompliance checks resource compliance with policies
func (a *Adapter) ValidateResourceCompliance(ctx context.Context, resourceType string, resource map[string]interface{}, rules []cloud.ValidationRule) ([]cloud.ValidationResult, error) {
	return []cloud.ValidationResult{}, nil
}

// DetectDrift compares planned state with actual cloud resources
func (a *Adapter) DetectDrift(ctx context.Context, plannedState map[string]interface{}, resourceType, resourceID string) (*cloud.ResourceStatus, error) {
	actualStatus, err := a.GetResourceStatus(ctx, resourceType, resourceID)
	if err != nil {
		return nil, err
	}

	if !actualStatus.Exists {
		actualStatus.DriftDetected = true
		actualStatus.DriftDetails = []string{"Resource does not exist in GCP"}
		return actualStatus, nil
	}

	driftDetails := []string{}

	// Check labels (GCP's version of tags)
	if plannedLabels, ok := plannedState["labels"].(map[string]interface{}); ok {
		for key, value := range plannedLabels {
			if actualValue, exists := actualStatus.Tags[key]; !exists {
				driftDetails = append(driftDetails, fmt.Sprintf("Label '%s' missing", key))
			} else if fmt.Sprint(value) != actualValue {
				driftDetails = append(driftDetails, fmt.Sprintf("Label '%s' differs: planned=%v, actual=%v", key, value, actualValue))
			}
		}
	}

	if len(driftDetails) > 0 {
		actualStatus.DriftDetected = true
		actualStatus.DriftDetails = driftDetails
	}

	return actualStatus, nil
}

// ListResources lists GCP resources of a given type
func (a *Adapter) ListResources(ctx context.Context, resourceType string) ([]string, error) {
	return nil, fmt.Errorf("listing not yet implemented for GCP")
}

// Close cleans up GCP adapter resources
func (a *Adapter) Close() error {
	if a.computeClient != nil {
		_ = a.computeClient.Close()
	}
	if a.storageClient != nil {
		_ = a.storageClient.Close()
	}
	return nil
}
