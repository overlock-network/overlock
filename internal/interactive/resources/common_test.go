package resources

import (
	"reflect"
	"testing"
	"time"

	condition "github.com/crossplane/crossplane-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name       string
		packageStr string
		expected   string
	}{
		{
			name:       "package with version",
			packageStr: "registry.io/provider-aws:v0.34.0",
			expected:   "v0.34.0",
		},
		{
			name:       "package with multiple colons",
			packageStr: "registry.io:8080/provider-aws:v0.34.0",
			expected:   "v0.34.0",
		},
		{
			name:       "package without version",
			packageStr: "registry.io/provider-aws",
			expected:   "latest",
		},
		{
			name:       "empty package string",
			packageStr: "",
			expected:   "latest",
		},
		{
			name:       "package with only colon",
			packageStr: "registry.io/provider-aws:",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractVersion(tt.packageStr)
			if result != tt.expected {
				t.Errorf("ExtractVersion(%q) = %q, want %q", tt.packageStr, result, tt.expected)
			}
		})
	}
}

func TestExtractInstallDate(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		expected  string
	}{
		{
			name:      "valid timestamp",
			timestamp: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			expected:  "2023-12-25",
		},
		{
			name:      "zero timestamp",
			timestamp: time.Time{},
			expected:  "Unknown",
		},
		{
			name:      "another valid timestamp",
			timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:  "2024-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractInstallDate(tt.timestamp)
			if result != tt.expected {
				t.Errorf("ExtractInstallDate(%v) = %q, want %q", tt.timestamp, result, tt.expected)
			}
		})
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name         string
		annotations  map[string]string
		packageStr   string
		resourceType string
		expected     string
	}{
		{
			name: "description annotation present",
			annotations: map[string]string{
				"description": "AWS Provider for Crossplane",
			},
			packageStr:   "registry.io/provider-aws:v0.34.0",
			resourceType: "provider",
			expected:     "AWS Provider for Crossplane",
		},
		{
			name: "meta.crossplane.io/description annotation present",
			annotations: map[string]string{
				"meta.crossplane.io/description": "Official AWS Provider",
			},
			packageStr:   "registry.io/provider-aws:v0.34.0",
			resourceType: "provider",
			expected:     "Official AWS Provider",
		},
		{
			name: "both annotations present - description takes precedence",
			annotations: map[string]string{
				"description":                    "First description",
				"meta.crossplane.io/description": "Second description",
			},
			packageStr:   "registry.io/provider-aws:v0.34.0",
			resourceType: "provider",
			expected:     "First description",
		},
		{
			name:         "no annotations, package with version",
			annotations:  nil,
			packageStr:   "registry.io/provider-aws:v0.34.0",
			resourceType: "provider",
			expected:     "Crossplane provider: provider-aws",
		},
		{
			name:         "no annotations, package without version",
			annotations:  map[string]string{},
			packageStr:   "registry.io/provider-aws",
			resourceType: "configuration",
			expected:     "Crossplane configuration: provider-aws",
		},
		{
			name:         "no annotations, complex package path",
			annotations:  nil,
			packageStr:   "registry.io/namespace/sub/provider-aws:v0.34.0",
			resourceType: "function",
			expected:     "Crossplane function: provider-aws",
		},
		{
			name:         "no annotations, empty package",
			annotations:  nil,
			packageStr:   "",
			resourceType: "provider",
			expected:     "Crossplane provider package",
		},
		{
			name:         "no annotations, only colon in package",
			annotations:  nil,
			packageStr:   ":",
			resourceType: "provider",
			expected:     "Crossplane provider: :",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractDescription(tt.annotations, tt.packageStr, tt.resourceType)
			if result != tt.expected {
				t.Errorf("ExtractDescription(%v, %q, %q) = %q, want %q",
					tt.annotations, tt.packageStr, tt.resourceType, result, tt.expected)
			}
		})
	}
}

func TestExtractStatusFromConditions_CrossplaneConditions(t *testing.T) {
	tests := []struct {
		name       string
		conditions []condition.Condition
		expected   string
	}{
		{
			name: "ready condition true",
			conditions: []condition.Condition{
				{
					Type:   condition.TypeReady,
					Status: corev1.ConditionTrue,
				},
			},
			expected: ConditionTypeReady,
		},
		{
			name: "healthy condition true, ready false",
			conditions: []condition.Condition{
				{
					Type:   condition.TypeReady,
					Status: corev1.ConditionFalse,
				},
				{
					Type:   "Healthy",
					Status: corev1.ConditionTrue,
				},
			},
			expected: ConditionTypeHealthy,
		},
		{
			name: "installed condition true",
			conditions: []condition.Condition{
				{
					Type:   "Installed",
					Status: corev1.ConditionTrue,
				},
			},
			expected: ConditionTypeInstalled,
		},
		{
			name: "synced condition true",
			conditions: []condition.Condition{
				{
					Type:   "Synced",
					Status: corev1.ConditionTrue,
				},
			},
			expected: ConditionTypeSynced,
		},
		{
			name: "failed condition true with reason",
			conditions: []condition.Condition{
				{
					Type:   "Failed",
					Status: corev1.ConditionTrue,
					Reason: "InstallationError",
				},
			},
			expected: "Failed: InstallationError",
		},
		{
			name: "failed condition true without reason",
			conditions: []condition.Condition{
				{
					Type:   "Failed",
					Status: corev1.ConditionTrue,
				},
			},
			expected: "Failed",
		},
		{
			name: "condition false with reason",
			conditions: []condition.Condition{
				{
					Type:   condition.TypeReady,
					Status: corev1.ConditionFalse,
					Reason: "PodNotReady",
				},
			},
			expected: "PodNotReady",
		},
		{
			name: "pending condition with message",
			conditions: []condition.Condition{
				{
					Type:    condition.TypeReady,
					Status:  corev1.ConditionFalse,
					Message: "Waiting for dependencies",
				},
			},
			expected: "Pending",
		},
		{
			name:       "empty conditions",
			conditions: []condition.Condition{},
			expected:   "Unknown",
		},
		{
			name:       "nil conditions",
			conditions: nil,
			expected:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractStatusFromConditions(tt.conditions)
			if result != tt.expected {
				t.Errorf("ExtractStatusFromConditions(%v) = %q, want %q", tt.conditions, result, tt.expected)
			}
		})
	}
}

func TestExtractStatusFromConditions_UnstructuredConditions(t *testing.T) {
	tests := []struct {
		name       string
		conditions []interface{}
		expected   string
	}{
		{
			name: "ready condition true",
			conditions: []interface{}{
				map[string]interface{}{
					"type":   "Ready",
					"status": "True",
				},
			},
			expected: ConditionTypeReady,
		},
		{
			name: "failed condition true with reason",
			conditions: []interface{}{
				map[string]interface{}{
					"type":   "Failed",
					"status": "True",
					"reason": "NetworkError",
				},
			},
			expected: "Failed: NetworkError",
		},
		{
			name: "multiple conditions with priority",
			conditions: []interface{}{
				map[string]interface{}{
					"type":   "Installed",
					"status": "True",
				},
				map[string]interface{}{
					"type":   "Ready",
					"status": "True",
				},
			},
			expected: ConditionTypeReady,
		},
		{
			name: "condition with invalid format",
			conditions: []interface{}{
				"invalid condition",
				map[string]interface{}{
					"type":   "Ready",
					"status": "True",
				},
			},
			expected: ConditionTypeReady,
		},
		{
			name:       "empty conditions",
			conditions: []interface{}{},
			expected:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractStatusFromConditions(tt.conditions)
			if result != tt.expected {
				t.Errorf("ExtractStatusFromConditions(%v) = %q, want %q", tt.conditions, result, tt.expected)
			}
		})
	}
}

func TestExtractStatusFromConditions_StructuredConditions(t *testing.T) {
	type TestCondition struct {
		Type    string
		Status  string
		Reason  string
		Message string
	}

	tests := []struct {
		name       string
		conditions interface{}
		expected   string
	}{
		{
			name: "ready condition true",
			conditions: []TestCondition{
				{
					Type:   "Ready",
					Status: "True",
				},
			},
			expected: ConditionTypeReady,
		},
		{
			name: "failed condition with reason",
			conditions: []TestCondition{
				{
					Type:   "Failed",
					Status: "True",
					Reason: "ConfigError",
				},
			},
			expected: "Failed: ConfigError",
		},
		{
			name: "multiple conditions with priority",
			conditions: []TestCondition{
				{
					Type:   "Synced",
					Status: "True",
				},
				{
					Type:   "Healthy",
					Status: "True",
				},
			},
			expected: ConditionTypeHealthy,
		},
		{
			name:       "non-slice input",
			conditions: "invalid",
			expected:   "Unknown",
		},
		{
			name:       "empty slice",
			conditions: []TestCondition{},
			expected:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractStatusFromConditions(tt.conditions)
			if result != tt.expected {
				t.Errorf("ExtractStatusFromConditions(%v) = %q, want %q", tt.conditions, result, tt.expected)
			}
		})
	}
}

func TestExtractStatusFromConditions_NilInput(t *testing.T) {
	result := ExtractStatusFromConditions(nil)
	expected := "Unknown"
	if result != expected {
		t.Errorf("ExtractStatusFromConditions(nil) = %q, want %q", result, expected)
	}
}

func TestGetFieldString(t *testing.T) {
	type TestStruct struct {
		Type    string
		Status  string
		Reason  string
		Message string
	}

	tests := []struct {
		name      string
		value     interface{}
		fieldName string
		expected  string
	}{
		{
			name: "valid struct with string field",
			value: TestStruct{
				Type:   "Ready",
				Status: "True",
			},
			fieldName: "Type",
			expected:  "Ready",
		},
		{
			name: "valid struct pointer with string field",
			value: &TestStruct{
				Status: "False",
			},
			fieldName: "Status",
			expected:  "False",
		},
		{
			name: "non-existent field",
			value: TestStruct{
				Type: "Ready",
			},
			fieldName: "NonExistent",
			expected:  "",
		},
		{
			name:      "non-struct input",
			value:     "not a struct",
			fieldName: "Type",
			expected:  "",
		},
		{
			name:      "nil pointer",
			value:     (*TestStruct)(nil),
			fieldName: "Type",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.value)
			result := getFieldString(v, tt.fieldName)
			if result != tt.expected {
				t.Errorf("getFieldString(%v, %q) = %q, want %q", tt.value, tt.fieldName, result, tt.expected)
			}
		})
	}
}
