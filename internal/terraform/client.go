// Package terraform provides utilities for interacting with Terraform.
package terraform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// Client wraps Terraform operations
type Client struct {
	workingDir   string
	terraformBin string
	backend      BackendConfig
	workspace    string
	envVars      map[string]string
}

// BackendConfig holds Terraform backend configuration
type BackendConfig struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

// PlanOutput represents the parsed output of terraform plan -json
type PlanOutput struct {
	FormatVersion    string                 `json:"format_version"`
	TerraformVersion string                 `json:"terraform_version"`
	PlannedValues    *StateValues           `json:"planned_values,omitempty"`
	ResourceChanges  []ResourceChange       `json:"resource_changes,omitempty"`
	Configuration    *Configuration         `json:"configuration,omitempty"`
	Variables        map[string]interface{} `json:"variables,omitempty"`
}

// StateValues represents Terraform state values
type StateValues struct {
	RootModule *Module `json:"root_module,omitempty"`
}

// Module represents a Terraform module
type Module struct {
	Resources    []Resource `json:"resources,omitempty"`
	ChildModules []Module   `json:"child_modules,omitempty"`
	Address      string     `json:"address,omitempty"`
}

// Resource represents a Terraform resource
type Resource struct {
	Address       string                 `json:"address"`
	Mode          string                 `json:"mode"` // "managed", "data"
	Type          string                 `json:"type"`
	Name          string                 `json:"name"`
	ProviderName  string                 `json:"provider_name"`
	SchemaVersion int                    `json:"schema_version"`
	Values        map[string]interface{} `json:"values"`
}

// ResourceChange represents a change to a resource
type ResourceChange struct {
	Address      string  `json:"address"`
	Mode         string  `json:"mode"`
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	ProviderName string  `json:"provider_name"`
	Change       *Change `json:"change"`
}

// Change represents the before/after values of a resource
type Change struct {
	Actions []string               `json:"actions"` // "create", "update", "delete", "no-op"
	Before  map[string]interface{} `json:"before"`
	After   map[string]interface{} `json:"after"`
}

// Configuration represents Terraform configuration
type Configuration struct {
	ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
	RootModule     *ConfigModule          `json:"root_module,omitempty"`
}

// ConfigModule represents module configuration
type ConfigModule struct {
	Resources   []ConfigResource       `json:"resources,omitempty"`
	ModuleCalls map[string]interface{} `json:"module_calls,omitempty"`
}

// ConfigResource represents a resource in configuration
type ConfigResource struct {
	Address      string                 `json:"address"`
	Mode         string                 `json:"mode"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	ProviderName string                 `json:"provider_config_key"`
	Expressions  map[string]interface{} `json:"expressions,omitempty"`
}

// NewClient creates a new Terraform client
func NewClient(workingDir string) (*Client, error) {
	if workingDir == "" {
		return nil, fmt.Errorf("working directory is required")
	}

	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("working directory does not exist: %s", workingDir)
	}

	// Find terraform binary
	terraformBin, err := exec.LookPath("terraform")
	if err != nil {
		return nil, fmt.Errorf("terraform binary not found in PATH: %w", err)
	}

	return &Client{
		workingDir:   workingDir,
		terraformBin: terraformBin,
		envVars:      make(map[string]string),
	}, nil
}

// SetEnvironment sets environment variables for Terraform execution
func (c *Client) SetEnvironment(key, value string) {
	c.envVars[key] = value
}

// SetWorkspace sets the Terraform workspace to use
func (c *Client) SetWorkspace(workspace string) {
	c.workspace = workspace
}

// Init runs terraform init
func (c *Client) Init(ctx context.Context, upgrade bool) error {
	args := []string{"init", "-no-color"}
	if upgrade {
		args = append(args, "-upgrade")
	}

	output, err := c.runCommand(ctx, args...)
	if err != nil {
		return fmt.Errorf("terraform init failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// Validate runs terraform validate
func (c *Client) Validate(ctx context.Context) error {
	output, err := c.runCommand(ctx, "validate", "-json")
	if err != nil {
		return fmt.Errorf("terraform validate failed: %w\nOutput: %s", err, output)
	}

	var result struct {
		Valid        bool `json:"valid"`
		ErrorCount   int  `json:"error_count"`
		WarningCount int  `json:"warning_count"`
		Diagnostics  []struct {
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
			Detail   string `json:"detail"`
		} `json:"diagnostics"`
	}

	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return fmt.Errorf("failed to parse validate output: %w", err)
	}

	if !result.Valid {
		var errMsgs []string
		for _, diag := range result.Diagnostics {
			if diag.Severity == "error" {
				errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", diag.Summary, diag.Detail))
			}
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, "; "))
	}

	return nil
}

// Plan runs terraform plan and returns the plan file path
func (c *Client) Plan(ctx context.Context, planFile string) error {
	args := []string{"plan", "-no-color", "-out=" + planFile}

	output, err := c.runCommand(ctx, args...)
	if err != nil {
		return fmt.Errorf("terraform plan failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// ShowJSON runs terraform show -json on a plan file
func (c *Client) ShowJSON(ctx context.Context, planFile string) (*PlanOutput, error) {
	output, err := c.runCommand(ctx, "show", "-json", planFile)
	if err != nil {
		return nil, fmt.Errorf("terraform show failed: %w\nOutput: %s", err, output)
	}

	var plan PlanOutput
	if err := json.Unmarshal([]byte(output), &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan output: %w", err)
	}

	return &plan, nil
}

// Apply runs terraform apply
func (c *Client) Apply(ctx context.Context, planFile string) error {
	args := []string{"apply", "-no-color", "-auto-approve"}
	if planFile != "" {
		args = append(args, planFile)
	}

	output, err := c.runCommand(ctx, args...)
	if err != nil {
		return fmt.Errorf("terraform apply failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// Destroy runs terraform destroy
func (c *Client) Destroy(ctx context.Context, autoApprove bool) error {
	args := []string{"destroy", "-no-color"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}

	output, err := c.runCommand(ctx, args...)
	if err != nil {
		return fmt.Errorf("terraform destroy failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// WorkspaceSelect selects a Terraform workspace
func (c *Client) WorkspaceSelect(ctx context.Context, workspace string) error {
	output, err := c.runCommand(ctx, "workspace", "select", workspace)
	if err != nil {
		// Try to create it if it doesn't exist
		output, err = c.runCommand(ctx, "workspace", "new", workspace)
		if err != nil {
			return fmt.Errorf("failed to create workspace %s: %w\nOutput: %s", workspace, err, output)
		}
	}

	c.workspace = workspace
	return nil
}

// GetProvider detects the cloud provider from Terraform configuration
func (c *Client) GetProvider(ctx context.Context) (string, error) {
	// Read all .tf files in the working directory
	files, err := filepath.Glob(filepath.Join(c.workingDir, "*.tf"))
	if err != nil {
		return "", fmt.Errorf("failed to list .tf files: %w", err)
	}

	providers := make(map[string]int)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		text := string(content)
		// Simple heuristic: count provider mentions
		if strings.Contains(text, `"aws"`) || strings.Contains(text, "provider \"aws\"") {
			providers["aws"]++
		}
		if strings.Contains(text, `"azurerm"`) || strings.Contains(text, "provider \"azurerm\"") {
			providers["azure"]++
		}
		if strings.Contains(text, `"google"`) || strings.Contains(text, "provider \"google\"") {
			providers["gcp"]++
		}
	}

	// Return the most common provider
	maxCount := 0
	detectedProvider := ""
	for provider, count := range providers {
		if count > maxCount {
			maxCount = count
			detectedProvider = provider
		}
	}

	if detectedProvider == "" {
		return "", fmt.Errorf("no cloud provider detected in Terraform configuration")
	}

	return detectedProvider, nil
}

// runCommand executes a Terraform command
func (c *Client) runCommand(ctx context.Context, args ...string) (string, error) {
	// Use direct execution - let the operating system handle path resolution
	cmd := exec.CommandContext(ctx, c.terraformBin, args...)
	cmd.Dir = c.workingDir

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range c.envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// On Windows, configure the subprocess creation attributes
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    false,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW to reduce visibility issues
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}

	if err != nil {
		return output, fmt.Errorf("command failed: %w", err)
	}

	return output, nil
}
