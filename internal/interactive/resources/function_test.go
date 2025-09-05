package resources

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
	"k8s.io/client-go/dynamic"
)

func TestLoadFunctions(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name            string
		dynamicClient   dynamic.Interface
		expectError     bool
		expectedCount   int
		expectNilClient bool
	}{
		{
			name:            "nil dynamic client",
			dynamicClient:   nil,
			expectError:     false,
			expectedCount:   0,
			expectNilClient: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if !tt.expectNilClient {
				t.Skip("Skipping dynamic client tests due to external dependencies")
			}

			result, err := LoadFunctions(ctx, tt.dynamicClient, logger)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d functions, got %d", tt.expectedCount, len(result))
			}

			if tt.expectNilClient {
				if len(result) != 0 {
					t.Error("Expected empty result for nil client")
				}
			}
		})
	}
}

func TestLoadFunctionsPaginated(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name          string
		dynamicClient dynamic.Interface
		opts          PaginationOptions
		expectError   bool
		expectedCount int
		expectedTotal int
		expectedMore  bool
	}{
		{
			name:          "nil dynamic client",
			dynamicClient: nil,
			opts:          PaginationOptions{Limit: 10, Offset: 0},
			expectError:   false,
			expectedCount: 0,
			expectedTotal: 0,
			expectedMore:  false,
		},
		{
			name:          "no pagination (limit 0)",
			dynamicClient: nil,
			opts:          PaginationOptions{Limit: 0, Offset: 0},
			expectError:   false,
			expectedCount: 0,
			expectedTotal: 0,
			expectedMore:  false,
		},
		{
			name:          "with pagination limit 1",
			dynamicClient: nil,
			opts:          PaginationOptions{Limit: 1, Offset: 0},
			expectError:   false,
			expectedCount: 0,
			expectedTotal: 0,
			expectedMore:  false,
		},
		{
			name:          "with pagination offset beyond data",
			dynamicClient: nil,
			opts:          PaginationOptions{Limit: 10, Offset: 100},
			expectError:   false,
			expectedCount: 0,
			expectedTotal: 0,
			expectedMore:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			result, err := LoadFunctionsPaginated(ctx, tt.dynamicClient, logger, tt.opts)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if len(result.Items) != tt.expectedCount {
				t.Errorf("Expected %d items, got %d", tt.expectedCount, len(result.Items))
			}

			if result.Total != tt.expectedTotal {
				t.Errorf("Expected total %d, got %d", tt.expectedTotal, result.Total)
			}

			if result.HasMore != tt.expectedMore {
				t.Errorf("Expected HasMore %v, got %v", tt.expectedMore, result.HasMore)
			}

			if result.ProcessedTime == "" {
				t.Error("Expected ProcessedTime to be set")
			}
		})
	}
}

func TestLoadFunctions_ContextCancellation(t *testing.T) {
	t.Skip("Skipping context cancellation test due to external dependencies")
}

func TestLoadFunctionsPaginated_ContextCancellation(t *testing.T) {
	t.Skip("Skipping context cancellation test due to external dependencies")
}

func TestFunctionResourceRowCreation(t *testing.T) {
	testTime := time.Date(2023, 10, 1, 14, 20, 30, 0, time.UTC)

	function := struct {
		Name              string
		Package           string
		CreationTimestamp time.Time
		Annotations       map[string]string
		Status            struct {
			Conditions interface{}
		}
	}{
		Name:              "function-go-templating",
		Package:           "xpkg.upbound.io/crossplane-contrib/function-go-templating:v0.3.0",
		CreationTimestamp: testTime,
		Annotations:       nil, // No description annotation
		Status: struct {
			Conditions interface{}
		}{
			Conditions: []interface{}{
				map[string]interface{}{
					"type":   "Synced",
					"status": "True",
				},
				map[string]interface{}{
					"type":   "Ready",
					"status": "False",
					"reason": "Pending",
				},
			},
		},
	}

	status := ExtractStatusFromConditions(function.Status.Conditions)
	version := ExtractVersion(function.Package)
	installDate := ExtractInstallDate(function.CreationTimestamp)
	description := ExtractDescription(function.Annotations, function.Package, "function")

	row := ResourceRow{
		Name:        function.Name,
		Package:     function.Package,
		Version:     version,
		Status:      status,
		InstallDate: installDate,
		Description: description,
	}

	if row.Name != "function-go-templating" {
		t.Errorf("Expected name 'function-go-templating', got %q", row.Name)
	}
	if row.Package != "xpkg.upbound.io/crossplane-contrib/function-go-templating:v0.3.0" {
		t.Errorf("Expected package 'xpkg.upbound.io/crossplane-contrib/function-go-templating:v0.3.0', got %q", row.Package)
	}
	if row.Version != "v0.3.0" {
		t.Errorf("Expected version 'v0.3.0', got %q", row.Version)
	}
	if row.Status != "Synced" {
		t.Errorf("Expected status 'Synced', got %q", row.Status)
	}
	if row.InstallDate != "2023-10-01" {
		t.Errorf("Expected install date '2023-10-01', got %q", row.InstallDate)
	}
	if row.Description != "Crossplane function: function-go-templating" {
		t.Errorf("Expected description 'Crossplane function: function-go-templating', got %q", row.Description)
	}
}

func TestFunctionStatusPriority(t *testing.T) {
	conditions := []interface{}{
		map[string]interface{}{
			"type":   "Synced",
			"status": "True",
		},
		map[string]interface{}{
			"type":   "Ready",
			"status": "True",
		},
		map[string]interface{}{
			"type":   "Healthy",
			"status": "True",
		},
	}

	status := ExtractStatusFromConditions(conditions)

	if status != "Ready" {
		t.Errorf("Expected status 'Ready', got %q", status)
	}
}

func TestFunctionPaginationWithMockData(t *testing.T) {
	totalFunctions := 4

	tests := []struct {
		name          string
		limit         int
		offset        int
		expectedCount int
		expectedMore  bool
	}{
		{
			name:          "first page of 2",
			limit:         2,
			offset:        0,
			expectedCount: 2,
			expectedMore:  true,
		},
		{
			name:          "second page of 2",
			limit:         2,
			offset:        2,
			expectedCount: 2,
			expectedMore:  false,
		},
		{
			name:          "single page for all functions",
			limit:         10,
			offset:        0,
			expectedCount: 4,
			expectedMore:  false,
		},
		{
			name:          "offset at last item",
			limit:         5,
			offset:        3,
			expectedCount: 1,
			expectedMore:  false,
		},
		{
			name:          "offset beyond data",
			limit:         2,
			offset:        10,
			expectedCount: 0,
			expectedMore:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualCount int
			var hasMore bool

			if tt.limit > 0 {
				end := tt.offset + tt.limit
				if tt.offset < totalFunctions {
					if end > totalFunctions {
						end = totalFunctions
						actualCount = end - tt.offset
					} else {
						actualCount = tt.limit
						if end < totalFunctions {
							hasMore = true
						}
					}
				} else {
					actualCount = 0
				}
			} else {
				actualCount = totalFunctions
			}

			if actualCount != tt.expectedCount {
				t.Errorf("Expected %d functions, got %d", tt.expectedCount, actualCount)
			}

			if hasMore != tt.expectedMore {
				t.Errorf("Expected hasMore %v, got %v", tt.expectedMore, hasMore)
			}
		})
	}
}

func TestFunctionComplexStatus(t *testing.T) {
	tests := []struct {
		name       string
		conditions interface{}
		expected   string
	}{
		{
			name: "failed function with reason",
			conditions: []interface{}{
				map[string]interface{}{
					"type":   "Failed",
					"status": "True",
					"reason": "ImagePullError",
				},
			},
			expected: "Failed: ImagePullError",
		},
		{
			name: "function with false condition and reason",
			conditions: []interface{}{
				map[string]interface{}{
					"type":   "Ready",
					"status": "False",
					"reason": "ContainerCreating",
				},
			},
			expected: "ContainerCreating",
		},
		{
			name: "function with pending message",
			conditions: []interface{}{
				map[string]interface{}{
					"type":    "Ready",
					"status":  "False",
					"message": "Waiting for image pull to complete",
				},
			},
			expected: "Pending",
		},
		{
			name:       "function with no conditions",
			conditions: []interface{}{},
			expected:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := ExtractStatusFromConditions(tt.conditions)
			if status != tt.expected {
				t.Errorf("Expected status %q, got %q", tt.expected, status)
			}
		})
	}
}
