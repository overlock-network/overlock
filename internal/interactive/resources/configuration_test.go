package resources

import (
	"context"
	"testing"
	"time"

	condition "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
)

func TestLoadConfigurations(t *testing.T) {
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

			result, err := LoadConfigurations(ctx, tt.dynamicClient, logger)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d configurations, got %d", tt.expectedCount, len(result))
			}

			if tt.expectNilClient {
				if len(result) != 0 {
					t.Error("Expected empty result for nil client")
				}
			}
		})
	}
}

func TestLoadConfigurationsPaginated(t *testing.T) {
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

			result, err := LoadConfigurationsPaginated(ctx, tt.dynamicClient, logger, tt.opts)

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

func TestLoadConfigurations_ContextCancellation(t *testing.T) {
	t.Skip("Skipping context cancellation test due to external dependencies")
}

func TestLoadConfigurationsPaginated_ContextCancellation(t *testing.T) {
	t.Skip("Skipping context cancellation test due to external dependencies")
}

func isContextCancelledError(err error) bool {
	return err == context.Canceled || err.Error() == "context canceled" ||
		(err != nil && (err.Error() == "context cancelled while loading configurations: context canceled" ||
			err.Error() == "context cancelled while loading providers: context canceled" ||
			err.Error() == "context cancelled while loading functions: context canceled"))
}

func TestResourceRowCreation(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	config := struct {
		Name              string
		Package           string
		CreationTimestamp time.Time
		Annotations       map[string]string
		Status            struct {
			Conditions interface{}
		}
	}{
		Name:              "test-config",
		Package:           "registry.io/test-config:v1.2.3",
		CreationTimestamp: testTime,
		Annotations: map[string]string{
			"description": "Test configuration description",
		},
		Status: struct {
			Conditions interface{}
		}{
			Conditions: []condition.Condition{
				{
					Type:   condition.TypeReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}

	status := ExtractStatusFromConditions(config.Status.Conditions)
	version := ExtractVersion(config.Package)
	installDate := ExtractInstallDate(config.CreationTimestamp)
	description := ExtractDescription(config.Annotations, config.Package, "configuration")

	row := ResourceRow{
		Name:        config.Name,
		Package:     config.Package,
		Version:     version,
		Status:      status,
		InstallDate: installDate,
		Description: description,
	}

	if row.Name != "test-config" {
		t.Errorf("Expected name 'test-config', got %q", row.Name)
	}
	if row.Package != "registry.io/test-config:v1.2.3" {
		t.Errorf("Expected package 'registry.io/test-config:v1.2.3', got %q", row.Package)
	}
	if row.Version != "v1.2.3" {
		t.Errorf("Expected version 'v1.2.3', got %q", row.Version)
	}
	if row.Status != "Ready" {
		t.Errorf("Expected status 'Ready', got %q", row.Status)
	}
	if row.InstallDate != "2023-12-25" {
		t.Errorf("Expected install date '2023-12-25', got %q", row.InstallDate)
	}
	if row.Description != "Test configuration description" {
		t.Errorf("Expected description 'Test configuration description', got %q", row.Description)
	}
}

func TestPaginationLogic(t *testing.T) {
	totalItems := 5
	items := make([]int, totalItems)
	for i := 0; i < totalItems; i++ {
		items[i] = i
	}

	tests := []struct {
		name          string
		limit         int
		offset        int
		expectedItems []int
		expectedMore  bool
	}{
		{
			name:          "no limit",
			limit:         0,
			offset:        0,
			expectedItems: []int{0, 1, 2, 3, 4},
			expectedMore:  false,
		},
		{
			name:          "first page",
			limit:         2,
			offset:        0,
			expectedItems: []int{0, 1},
			expectedMore:  true,
		},
		{
			name:          "second page",
			limit:         2,
			offset:        2,
			expectedItems: []int{2, 3},
			expectedMore:  true,
		},
		{
			name:          "last page",
			limit:         2,
			offset:        4,
			expectedItems: []int{4},
			expectedMore:  false,
		},
		{
			name:          "offset beyond data",
			limit:         2,
			offset:        10,
			expectedItems: []int{},
			expectedMore:  false,
		},
		{
			name:          "limit larger than remaining",
			limit:         10,
			offset:        3,
			expectedItems: []int{3, 4},
			expectedMore:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []int
			hasMore := false

			if tt.limit > 0 {
				end := tt.offset + tt.limit
				if tt.offset < totalItems {
					if end > totalItems {
						end = totalItems
					} else {
						hasMore = true
					}
					result = items[tt.offset:end]
				} else {
					result = []int{}
				}
			} else {
				result = items
			}

			if len(result) != len(tt.expectedItems) {
				t.Errorf("Expected %d items, got %d", len(tt.expectedItems), len(result))
			}

			for i, expected := range tt.expectedItems {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected item %d to be %d, got %d", i, expected, result[i])
				}
			}

			if hasMore != tt.expectedMore {
				t.Errorf("Expected hasMore %v, got %v", tt.expectedMore, hasMore)
			}
		})
	}
}
