package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vijayaxai/terraship/internal/cloud"
)

func TestRulesEngine_RequiredTags(t *testing.T) {
	engine := &Engine{
		policy: &Policy{
			Rules: []cloud.ValidationRule{
				{
					Name:          "required-tags",
					Severity:      "error",
					Enabled:       true,
					ResourceTypes: []string{"aws_*"},
					Conditions: map[string]interface{}{
						"tags.required": []interface{}{"Environment", "Owner"},
					},
					Message: "Required tags missing",
				},
			},
		},
	}

	tests := []struct {
		name       string
		resource   map[string]interface{}
		shouldPass bool
	}{
		{
			name: "all tags present",
			resource: map[string]interface{}{
				"tags": map[string]interface{}{
					"Environment": "prod",
					"Owner":       "team@example.com",
				},
			},
			shouldPass: true,
		},
		{
			name: "missing tags",
			resource: map[string]interface{}{
				"tags": map[string]interface{}{
					"Environment": "prod",
				},
			},
			shouldPass: false,
		},
		{
			name:       "no tags",
			resource:   map[string]interface{}{},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := engine.GetRulesForResource("aws_instance")
			assert.Len(t, rules, 1)

			result := engine.EvaluateRule(rules[0], tt.resource)
			assert.Equal(t, tt.shouldPass, result.Passed)
		})
	}
}

func TestRulesEngine_Encryption(t *testing.T) {
	engine := &Engine{
		policy: &Policy{
			Rules: []cloud.ValidationRule{
				{
					Name:          "encryption-enabled",
					Severity:      "error",
					Enabled:       true,
					ResourceTypes: []string{"aws_s3_bucket"},
					Conditions: map[string]interface{}{
						"encryption.enabled": true,
					},
					Message: "Encryption must be enabled",
				},
			},
		},
	}

	tests := []struct {
		name       string
		resource   map[string]interface{}
		shouldPass bool
	}{
		{
			name: "encryption enabled",
			resource: map[string]interface{}{
				"encrypted": true,
			},
			shouldPass: true,
		},
		{
			name: "encryption disabled",
			resource: map[string]interface{}{
				"encrypted": false,
			},
			shouldPass: false,
		},
		{
			name:       "no encryption field",
			resource:   map[string]interface{}{},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := engine.GetRulesForResource("aws_s3_bucket")
			assert.Len(t, rules, 1)

			result := engine.EvaluateRule(rules[0], tt.resource)
			assert.Equal(t, tt.shouldPass, result.Passed, result.Details)
		})
	}
}
