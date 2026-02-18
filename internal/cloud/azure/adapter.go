// Package azure implements the Azure cloud adapter.
package azure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/vijayaxai/terraship/internal/cloud"
)

// Adapter implements cloud.Adapter for Azure
type Adapter struct {
	cred            azcore.TokenCredential
	subscriptionID  string
	resourcesClient *armresources.Client
	computeClient   *armcompute.VirtualMachinesClient
	storageClient   *armstorage.AccountsClient
}

// NewAdapter creates a new Azure adapter
func NewAdapter() *Adapter {
	return &Adapter{}
}

// Name returns the provider name
func (a *Adapter) Name() cloud.Provider {
	return cloud.ProviderAzure
}

// Initialize sets up the Azure adapter
func (a *Adapter) Initialize(ctx context.Context, cloudConfig cloud.CloudConfig) error {
	var err error

	// Get subscription ID
	a.subscriptionID = cloudConfig.AzureSubscriptionID
	if a.subscriptionID == "" {
		a.subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
	if a.subscriptionID == "" {
		return fmt.Errorf("Azure subscription ID is required")
	}

	// Create credential
	if cloudConfig.AzureClientID != "" && cloudConfig.AzureClientSecret != "" && cloudConfig.AzureTenantID != "" {
		// Use client secret credential
		a.cred, err = azidentity.NewClientSecretCredential(
			cloudConfig.AzureTenantID,
			cloudConfig.AzureClientID,
			cloudConfig.AzureClientSecret,
			nil,
		)
	} else {
		// Use default credential chain (Azure CLI, managed identity, etc.)
		a.cred, err = azidentity.NewDefaultAzureCredential(nil)
	}

	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %w", err)
	}

	// Initialize clients
	a.resourcesClient, err = armresources.NewClient(a.subscriptionID, a.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create resources client: %w", err)
	}

	a.computeClient, err = armcompute.NewVirtualMachinesClient(a.subscriptionID, a.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create compute client: %w", err)
	}

	a.storageClient, err = armstorage.NewAccountsClient(a.subscriptionID, a.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}

	return nil
}

// DetectProvider attempts to detect if Azure is the provider
func (a *Adapter) DetectProvider(ctx context.Context) (bool, float64, error) {
	confidence := 0.0

	// Check for Azure environment variables
	if os.Getenv("AZURE_SUBSCRIPTION_ID") != "" {
		confidence += 0.3
	}
	if os.Getenv("AZURE_TENANT_ID") != "" {
		confidence += 0.2
	}
	if os.Getenv("AZURE_CLIENT_ID") != "" {
		confidence += 0.2
	}
	if os.Getenv("ARM_SUBSCRIPTION_ID") != "" {
		confidence += 0.3
	}

	return confidence > 0.5, confidence, nil
}

// ValidateCredentials checks if Azure credentials are valid
func (a *Adapter) ValidateCredentials(ctx context.Context) error {
	// Try to list resource groups as a lightweight validation
	pager := a.resourcesClient.NewListPager(nil)
	if !pager.More() {
		return nil
	}

	_, err := pager.NextPage(ctx)
	if err != nil {
		return fmt.Errorf("Azure credentials validation failed: %w", err)
	}

	return nil
}

// GetResourceStatus retrieves the current status of an Azure resource
func (a *Adapter) GetResourceStatus(ctx context.Context, resourceType, resourceID string) (*cloud.ResourceStatus, error) {
	status := &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Properties:   make(map[string]interface{}),
	}

	switch {
	case strings.HasPrefix(resourceType, "azurerm_virtual_machine"):
		return a.getVMStatus(ctx, resourceID)
	case strings.HasPrefix(resourceType, "azurerm_storage_account"):
		return a.getStorageAccountStatus(ctx, resourceID)
	case strings.HasPrefix(resourceType, "azurerm_resource_group"):
		return a.getResourceGroupStatus(ctx, resourceID)
	default:
		return status, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (a *Adapter) getVMStatus(ctx context.Context, resourceID string) (*cloud.ResourceStatus, error) {
	// Parse resource ID: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Compute/virtualMachines/{name}
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		return nil, fmt.Errorf("invalid Azure resource ID format")
	}

	resourceGroup := parts[4]
	vmName := parts[8]

	vm, err := a.computeClient.Get(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   resourceID,
			ResourceType: "azurerm_virtual_machine",
			Exists:       false,
		}, nil
	}

	tags := make(map[string]string)
	if vm.Tags != nil {
		for key, value := range vm.Tags {
			if value != nil {
				tags[key] = *value
			}
		}
	}

	status := &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: "azurerm_virtual_machine",
		Exists:       true,
		Tags:         tags,
		Properties:   make(map[string]interface{}),
	}

	if vm.Properties != nil {
		if vm.Properties.HardwareProfile != nil {
			status.Properties["vm_size"] = vm.Properties.HardwareProfile.VMSize
		}
		if vm.Properties.ProvisioningState != nil {
			status.State = *vm.Properties.ProvisioningState
		}
	}

	return status, nil
}

func (a *Adapter) getStorageAccountStatus(ctx context.Context, resourceID string) (*cloud.ResourceStatus, error) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		return nil, fmt.Errorf("invalid Azure resource ID format")
	}

	resourceGroup := parts[4]
	accountName := parts[8]

	account, err := a.storageClient.GetProperties(ctx, resourceGroup, accountName, nil)
	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   resourceID,
			ResourceType: "azurerm_storage_account",
			Exists:       false,
		}, nil
	}

	tags := make(map[string]string)
	if account.Tags != nil {
		for key, value := range account.Tags {
			if value != nil {
				tags[key] = *value
			}
		}
	}

	status := &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: "azurerm_storage_account",
		Exists:       true,
		Tags:         tags,
		Properties:   make(map[string]interface{}),
	}

	if account.Properties != nil {
		if account.Properties.Encryption != nil {
			status.Properties["encryption_enabled"] = true
		}
		if account.Properties.ProvisioningState != nil {
			status.State = string(*account.Properties.ProvisioningState)
		}
	}

	return status, nil
}

func (a *Adapter) getResourceGroupStatus(ctx context.Context, resourceID string) (*cloud.ResourceStatus, error) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		// Assume it's a resource group if path is short
		if len(parts) >= 5 {
			// Format: /subscriptions/{subscription}/resourceGroups/{rg}
			return &cloud.ResourceStatus{
				ResourceID:   resourceID,
				ResourceType: "azurerm_resource_group",
				Exists:       true, // Assume exists for now
			}, nil
		}
		return nil, fmt.Errorf("invalid Azure resource ID format")
	}

	resourceProvider := parts[6]
	resourceType := parts[7]
	resourceName := parts[8]
	resourceGroup := parts[4]

	// Get the resource
	res, err := a.resourcesClient.GetByID(ctx, resourceID, "2021-04-01", nil)
	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   resourceID,
			ResourceType: fmt.Sprintf("%s/%s", resourceProvider, resourceType),
			Exists:       false,
		}, nil
	}

	tags := make(map[string]string)
	if res.Tags != nil {
		for key, value := range res.Tags {
			if value != nil {
				tags[key] = *value
			}
		}
	}

	return &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: fmt.Sprintf("%s/%s", resourceProvider, resourceType),
		Exists:       true,
		Tags:         tags,
		Properties: map[string]interface{}{
			"location":      res.Location,
			"resourceGroup": resourceGroup,
			"name":          resourceName,
		},
	}, nil
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
		actualStatus.DriftDetails = []string{"Resource does not exist in Azure"}
		return actualStatus, nil
	}

	driftDetails := []string{}

	// Check tags
	if plannedTags, ok := plannedState["tags"].(map[string]interface{}); ok {
		for key, value := range plannedTags {
			if actualValue, exists := actualStatus.Tags[key]; !exists {
				driftDetails = append(driftDetails, fmt.Sprintf("Tag '%s' missing", key))
			} else if fmt.Sprint(value) != actualValue {
				driftDetails = append(driftDetails, fmt.Sprintf("Tag '%s' differs: planned=%v, actual=%v", key, value, actualValue))
			}
		}
	}

	if len(driftDetails) > 0 {
		actualStatus.DriftDetected = true
		actualStatus.DriftDetails = driftDetails
	}

	return actualStatus, nil
}

// ListResources lists Azure resources of a given type
func (a *Adapter) ListResources(ctx context.Context, resourceType string) ([]string, error) {
	return nil, fmt.Errorf("listing not yet implemented for Azure")
}

// Close cleans up Azure adapter resources
func (a *Adapter) Close() error {
	return nil
}
