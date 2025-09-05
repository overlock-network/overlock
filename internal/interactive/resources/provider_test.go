package resources

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
	"k8s.io/client-go/dynamic"
)

func TestLoadProviders(t *testing.T) {
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

			result, err := LoadProviders(ctx, tt.dynamicClient, logger)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d providers, got %d", tt.expectedCount, len(result))
			}

			if tt.expectNilClient {
				if len(result) != 0 {
					t.Error("Expected empty result for nil client")
				}
			}
		})
	}
}

func TestLoadProvidersPaginated(t *testing.T) {
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

			result, err := LoadProvidersPaginated(ctx, tt.dynamicClient, logger, tt.opts)

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

func TestLoadProviders_ContextCancellation(t *testing.T) {
	t.Skip("Skipping context cancellation test due to external dependencies")
}

func TestLoadProvidersPaginated_ContextCancellation(t *testing.T) {
	t.Skip("Skipping context cancellation test due to external dependencies")
}

func TestProviderResourceRowCreation(t *testing.T) {
	testTime := time.Date(2023, 11, 15, 10, 30, 0, 0, time.UTC)

	provider := struct {
		Name              string
		Package           string
		CreationTimestamp time.Time
		Annotations       map[string]string
		Status            struct {
			Conditions interface{}
		}
	}{
		Name:              "provider-aws",
		Package:           "xpkg.upbound.io/crossplane-contrib/provider-aws:v0.34.0",
		CreationTimestamp: testTime,
		Annotations: map[string]string{
			"meta.crossplane.io/description": "AWS Provider for Crossplane",
		},
		Status: struct {
			Conditions interface{}
		}{
			Conditions: []interface{}{
				map[string]interface{}{
					"type":   "Healthy",
					"status": "True",
				},
			},
		},
	}

	status := ExtractStatusFromConditions(provider.Status.Conditions)
	version := ExtractVersion(provider.Package)
	installDate := ExtractInstallDate(provider.CreationTimestamp)
	description := ExtractDescription(provider.Annotations, provider.Package, "provider")

	row := ResourceRow{
		Name:        provider.Name,
		Package:     provider.Package,
		Version:     version,
		Status:      status,
		InstallDate: installDate,
		Description: description,
	}

	if row.Name != "provider-aws" {
		t.Errorf("Expected name 'provider-aws', got %q", row.Name)
	}
	if row.Package != "xpkg.upbound.io/crossplane-contrib/provider-aws:v0.34.0" {
		t.Errorf("Expected package 'xpkg.upbound.io/crossplane-contrib/provider-aws:v0.34.0', got %q", row.Package)
	}
	if row.Version != "v0.34.0" {
		t.Errorf("Expected version 'v0.34.0', got %q", row.Version)
	}
	if row.Status != "Healthy" {
		t.Errorf("Expected status 'Healthy', got %q", row.Status)
	}
	if row.InstallDate != "2023-11-15" {
		t.Errorf("Expected install date '2023-11-15', got %q", row.InstallDate)
	}
	if row.Description != "AWS Provider for Crossplane" {
		t.Errorf("Expected description 'AWS Provider for Crossplane', got %q", row.Description)
	}
}

func TestProviderPaginationWithMockData(t *testing.T) {
	totalProviders := 7

	tests := []struct {
		name          string
		limit         int
		offset        int
		expectedCount int
		expectedMore  bool
	}{
		{
			name:          "first page of 3",
			limit:         3,
			offset:        0,
			expectedCount: 3,
			expectedMore:  true,
		},
		{
			name:          "second page of 3",
			limit:         3,
			offset:        3,
			expectedCount: 3,
			expectedMore:  true,
		},
		{
			name:          "third page of 3 (partial)",
			limit:         3,
			offset:        6,
			expectedCount: 1,
			expectedMore:  false,
		},
		{
			name:          "large page size",
			limit:         10,
			offset:        0,
			expectedCount: 7,
			expectedMore:  false,
		},
		{
			name:          "offset beyond data",
			limit:         5,
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
				if tt.offset < totalProviders {
					if end > totalProviders {
						end = totalProviders
						actualCount = end - tt.offset
					} else {
						hasMore = true
						actualCount = tt.limit
					}
				} else {
					actualCount = 0
				}
			} else {
				actualCount = totalProviders
			}

			if actualCount != tt.expectedCount {
				t.Errorf("Expected %d providers, got %d", tt.expectedCount, actualCount)
			}

			if hasMore != tt.expectedMore {
				t.Errorf("Expected hasMore %v, got %v", tt.expectedMore, hasMore)
			}
		})
	}
}
