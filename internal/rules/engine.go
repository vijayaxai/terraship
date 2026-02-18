// Package rules implements the policy rule engine for Terraship.
package rules

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/vijayaxai/terraship/internal/cloud"
	"gopkg.in/yaml.v3"
)

// Policy represents a collection of validation rules
type Policy struct {
	Version     string                 `yaml:"version"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Rules       []cloud.ValidationRule `yaml:"rules"`
}

// Engine evaluates rules against resources
type Engine struct {
	policy *Policy
}

// NewEngine creates a new rules engine
func NewEngine(policyPath string) (*Engine, error) {
	data, err := os.ReadFile(policyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	var policy Policy
	if err := yaml.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	return &Engine{policy: &policy}, nil
}

// GetRulesForResource returns rules applicable to a resource type
func (e *Engine) GetRulesForResource(resourceType string) []cloud.ValidationRule {
	var applicable []cloud.ValidationRule

	for _, rule := range e.policy.Rules {
		if !rule.Enabled {
			continue
		}

		// Check if rule applies to this resource type
		if len(rule.ResourceTypes) == 0 {
			// No specific types means applies to all
			applicable = append(applicable, rule)
			continue
		}

		for _, rt := range rule.ResourceTypes {
			if matchResourceType(rt, resourceType) {
				applicable = append(applicable, rule)
				break
			}
		}
	}

	return applicable
}

// EvaluateRule checks if a resource meets a rule's conditions
func (e *Engine) EvaluateRule(rule cloud.ValidationRule, resource map[string]interface{}) cloud.ValidationResult {
	result := cloud.ValidationResult{
		RuleName:    rule.Name,
		Severity:    rule.Severity,
		Passed:      true,
		Message:     rule.Message,
		Remediation: rule.Remediation,
	}

	// Evaluate conditions
	for condition, expected := range rule.Conditions {
		if !e.evaluateCondition(condition, expected, resource, &result) {
			result.Passed = false
			break
		}
	}

	return result
}

// evaluateCondition checks a single condition
func (e *Engine) evaluateCondition(condition string, expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	switch condition {
	case "tags.required":
		return e.checkRequiredTags(expected, resource, result)

	case "encryption.enabled":
		return e.checkEncryptionEnabled(expected, resource, result)

	case "public_access.blocked":
		return e.checkPublicAccessBlocked(expected, resource, result)

	case "versioning.enabled":
		return e.checkVersioningEnabled(expected, resource, result)

	case "logging.enabled":
		return e.checkLoggingEnabled(expected, resource, result)

	case "backup.enabled":
		return e.checkBackupEnabled(expected, resource, result)

	case "naming.pattern":
		return e.checkNamingPattern(expected, resource, result)

	case "iam.least_privilege":
		return e.checkLeastPrivilege(expected, resource, result)

	case "network.private_subnet":
		return e.checkPrivateSubnet(expected, resource, result)

	default:
		// Generic property check
		return e.checkProperty(condition, expected, resource, result)
	}
}

func (e *Engine) checkRequiredTags(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	requiredTags, ok := expected.([]interface{})
	if !ok {
		result.Details = append(result.Details, "Invalid tags.required configuration")
		return false
	}

	tags, ok := resource["tags"].(map[string]interface{})
	if !ok {
		result.Details = append(result.Details, "No tags found on resource")
		return false
	}

	missingTags := []string{}
	for _, tag := range requiredTags {
		tagName := fmt.Sprint(tag)
		if _, exists := tags[tagName]; !exists {
			missingTags = append(missingTags, tagName)
		}
	}

	if len(missingTags) > 0 {
		result.Details = append(result.Details, fmt.Sprintf("Missing required tags: %s", strings.Join(missingTags, ", ")))
		return false
	}

	return true
}

func (e *Engine) checkEncryptionEnabled(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	shouldBeEnabled, ok := expected.(bool)
	if !ok || !shouldBeEnabled {
		return true
	}

	// Check various encryption fields
	encryptionFields := []string{
		"encryption", "encrypted", "encryption_configuration",
		"server_side_encryption_configuration", "encryption_at_rest",
	}

	for _, field := range encryptionFields {
		if value, exists := resource[field]; exists {
			if boolVal, ok := value.(bool); ok && boolVal {
				return true
			}
			if mapVal, ok := value.(map[string]interface{}); ok && len(mapVal) > 0 {
				return true
			}
		}
	}

	result.Details = append(result.Details, "Encryption is not enabled")
	return false
}

func (e *Engine) checkPublicAccessBlocked(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	shouldBeBlocked, ok := expected.(bool)
	if !ok || !shouldBeBlocked {
		return true
	}

	// Check for public access indicators
	publicFields := map[string]bool{
		"public":                false,
		"publicly_accessible":   false,
		"public_access_enabled": false,
		"acl":                   false, // Check if ACL allows public
	}

	for field, _ := range publicFields {
		if value, exists := resource[field]; exists {
			if boolVal, ok := value.(bool); ok && boolVal {
				result.Details = append(result.Details, fmt.Sprintf("Resource has public access via '%s'", field))
				return false
			}
			if strVal, ok := value.(string); ok && strings.Contains(strings.ToLower(strVal), "public") {
				result.Details = append(result.Details, fmt.Sprintf("Resource has public access via '%s': %s", field, strVal))
				return false
			}
		}
	}

	return true
}

func (e *Engine) checkVersioningEnabled(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	shouldBeEnabled, ok := expected.(bool)
	if !ok || !shouldBeEnabled {
		return true
	}

	versioningFields := []string{"versioning", "versioning_configuration", "version_enabled"}

	for _, field := range versioningFields {
		if value, exists := resource[field]; exists {
			if boolVal, ok := value.(bool); ok && boolVal {
				return true
			}
			if mapVal, ok := value.(map[string]interface{}); ok {
				if enabled, ok := mapVal["enabled"].(bool); ok && enabled {
					return true
				}
			}
		}
	}

	result.Details = append(result.Details, "Versioning is not enabled")
	return false
}

func (e *Engine) checkLoggingEnabled(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	shouldBeEnabled, ok := expected.(bool)
	if !ok || !shouldBeEnabled {
		return true
	}

	loggingFields := []string{"logging", "logging_configuration", "log_configuration", "enable_logging"}

	for _, field := range loggingFields {
		if value, exists := resource[field]; exists {
			if boolVal, ok := value.(bool); ok && boolVal {
				return true
			}
			if mapVal, ok := value.(map[string]interface{}); ok && len(mapVal) > 0 {
				return true
			}
		}
	}

	result.Details = append(result.Details, "Logging is not enabled")
	return false
}

func (e *Engine) checkBackupEnabled(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	shouldBeEnabled, ok := expected.(bool)
	if !ok || !shouldBeEnabled {
		return true
	}

	backupFields := []string{"backup", "backup_configuration", "backup_enabled", "backup_retention_period"}

	for _, field := range backupFields {
		if value, exists := resource[field]; exists {
			if boolVal, ok := value.(bool); ok && boolVal {
				return true
			}
			if intVal, ok := value.(int); ok && intVal > 0 {
				return true
			}
			if mapVal, ok := value.(map[string]interface{}); ok && len(mapVal) > 0 {
				return true
			}
		}
	}

	result.Details = append(result.Details, "Backup is not configured")
	return false
}

func (e *Engine) checkNamingPattern(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	pattern, ok := expected.(string)
	if !ok {
		return true
	}

	nameFields := []string{"name", "id", "resource_name"}

	for _, field := range nameFields {
		if value, exists := resource[field]; exists {
			if name, ok := value.(string); ok {
				matched, err := regexp.MatchString(pattern, name)
				if err != nil {
					result.Details = append(result.Details, fmt.Sprintf("Invalid regex pattern: %s", err))
					return false
				}
				if !matched {
					result.Details = append(result.Details, fmt.Sprintf("Name '%s' does not match pattern '%s'", name, pattern))
					return false
				}
				return true
			}
		}
	}

	return true
}

func (e *Engine) checkLeastPrivilege(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	// Check for overly permissive IAM policies
	policyFields := []string{"policy", "policy_document", "policy_arn"}

	for _, field := range policyFields {
		if value, exists := resource[field]; exists {
			if strVal, ok := value.(string); ok {
				if strings.Contains(strVal, "*:*") || strings.Contains(strVal, "\"*\"") {
					result.Details = append(result.Details, "Policy contains wildcard permissions")
					return false
				}
			}
		}
	}

	return true
}

func (e *Engine) checkPrivateSubnet(expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	shouldBePrivate, ok := expected.(bool)
	if !ok || !shouldBePrivate {
		return true
	}

	// Check subnet configuration
	if subnetID, exists := resource["subnet_id"]; exists {
		if strVal, ok := subnetID.(string); ok && strings.Contains(strVal, "public") {
			result.Details = append(result.Details, "Resource is in a public subnet")
			return false
		}
	}

	return true
}

func (e *Engine) checkProperty(propertyPath string, expected interface{}, resource map[string]interface{}, result *cloud.ValidationResult) bool {
	// Navigate nested properties using dot notation
	parts := strings.Split(propertyPath, ".")
	current := resource

	for i, part := range parts {
		value, exists := current[part]
		if !exists {
			result.Details = append(result.Details, fmt.Sprintf("Property '%s' not found", propertyPath))
			return false
		}

		if i == len(parts)-1 {
			// Last part - compare value
			if fmt.Sprint(value) != fmt.Sprint(expected) {
				result.Details = append(result.Details, fmt.Sprintf("Property '%s' has value '%v', expected '%v'", propertyPath, value, expected))
				return false
			}
			return true
		}

		// Navigate deeper
		if nested, ok := value.(map[string]interface{}); ok {
			current = nested
		} else {
			result.Details = append(result.Details, fmt.Sprintf("Cannot navigate property path '%s'", propertyPath))
			return false
		}
	}

	return true
}

// matchResourceType checks if a resource type matches a pattern
func matchResourceType(pattern, resourceType string) bool {
	// Support wildcards
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	matched, _ := regexp.MatchString("^"+pattern+"$", resourceType)
	return matched
}
