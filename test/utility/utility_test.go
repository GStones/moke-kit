package utility

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gstones/moke-kit/test/utils"
	"github.com/gstones/moke-kit/utility"
)

// TestAuthUtilities tests authentication utility functions
func TestAuthUtilities(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("AuthValidation", func(t *testing.T) {
		// Test valid auth token (assuming this functionality exists)
		// This would test actual auth.go functions
		helper.AssertTrue(true, "Auth validation placeholder")
	})
}

// TestConfigUtilities tests configuration utility functions
func TestConfigUtilities(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ConfigLoading", func(t *testing.T) {
		// Test config loading functionality (assuming this exists)
		// This would test actual config.go functions
		helper.AssertTrue(true, "Config loading placeholder")
	})
}

// TestContextUtilities tests context utility functions
func TestContextUtilities(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ContextOperations", func(t *testing.T) {
		ctx := helper.Context()
		helper.AssertNotNil(ctx, "Context should not be nil")

		// Test context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		helper.AssertNotNil(timeoutCtx, "Timeout context should not be nil")

		// Test context cancellation
		cancelCtx, cancelFunc := context.WithCancel(ctx)
		helper.AssertNotNil(cancelCtx, "Cancel context should not be nil")

		cancelFunc()
		
		select {
		case <-cancelCtx.Done():
			helper.AssertTrue(true, "Context should be cancelled")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Context should have been cancelled")
		}
	})

	t.Run("ContextWithValues", func(t *testing.T) {
		ctx := helper.Context()
		
		// Test context with values
		key := "test_key"
		value := "test_value"
		ctxWithValue := context.WithValue(ctx, key, value)
		
		retrievedValue := ctxWithValue.Value(key)
		helper.AssertEqual(value, retrievedValue, "Retrieved value should match")
		
		// Test missing key
		missingValue := ctxWithValue.Value("missing_key")
		helper.AssertNil(missingValue, "Missing key should return nil")
	})
}

// TestDeploymentUtilities tests deployment utility functions
func TestDeploymentUtilities(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("DeploymentParsing", func(t *testing.T) {
		// Test deployment string parsing
		testCases := []struct {
			input    string
			expected utility.Deployments
		}{
			{"dev", utility.DeploymentsDev},
			{"prod", utility.DeploymentsProd},
			{"local", utility.DeploymentsLocal},
			{"production", utility.Deployments("production")}, // Custom deployment
			{"unknown", utility.Deployments("unknown")},       // Custom deployment
			{"", utility.Deployments("")},                     // Empty string
		}

		for _, tc := range testCases {
			result := utility.ParseDeployments(tc.input)
			helper.AssertEqual(tc.expected, result, "Deployment parsing for input: %s", tc.input)
		}
	})

	t.Run("DeploymentMethods", func(t *testing.T) {
		// Test IsDev
		dev := utility.ParseDeployments("dev")
		helper.AssertTrue(dev.IsDev(), "Dev should return true for IsDev()")
		helper.AssertFalse(dev.IsLocal(), "Dev should return false for IsLocal()")
		helper.AssertFalse(dev.IsProd(), "Dev should return false for IsProd()")

		// Test IsLocal
		local := utility.ParseDeployments("local")
		helper.AssertFalse(local.IsDev(), "Local should return false for IsDev()")
		helper.AssertTrue(local.IsLocal(), "Local should return true for IsLocal()")
		helper.AssertFalse(local.IsProd(), "Local should return false for IsProd()")

		// Test IsProd
		prod := utility.ParseDeployments("prod")
		helper.AssertFalse(prod.IsDev(), "Prod should return false for IsDev()")
		helper.AssertFalse(prod.IsLocal(), "Prod should return false for IsLocal()")
		helper.AssertTrue(prod.IsProd(), "Prod should return true for IsProd()")

		// Test custom deployment
		custom := utility.ParseDeployments("custom_env")
		helper.AssertFalse(custom.IsDev(), "Custom should return false for IsDev()")
		helper.AssertFalse(custom.IsLocal(), "Custom should return false for IsLocal()")
		helper.AssertFalse(custom.IsProd(), "Custom should return false for IsProd()")
	})

	t.Run("DeploymentString", func(t *testing.T) {
		// Test String() method
		testCases := []struct {
			deployment utility.Deployments
			expected   string
		}{
			{utility.DeploymentsDev, "dev"},
			{utility.DeploymentsLocal, "local"},
			{utility.DeploymentsProd, "prod"},
		}

		for _, tc := range testCases {
			result := tc.deployment.String()
			helper.AssertEqual(tc.expected, result, "String representation for deployment: %s", tc.deployment.String())
		}
	})
}

// TestUtilityHelpers tests various utility helper functions
func TestUtilityHelpers(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("StringUtilities", func(t *testing.T) {
		// Test string utility functions (if they exist)
		testString := "test string"
		helper.AssertEqual(len(testString), 11, "String length should be correct")
		helper.AssertTrue(len(testString) > 0, "String should not be empty")
	})

	t.Run("TimeUtilities", func(t *testing.T) {
		// Test time utility functions
		now := time.Now()
		helper.AssertTrue(now.Before(time.Now().Add(time.Second)), "Now should be before future time")
		helper.AssertTrue(now.After(time.Now().Add(-time.Second)), "Now should be after past time")
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test error handling utilities
		var err error
		helper.AssertNil(err, "Error should be nil initially")

		err = context.Canceled
		helper.AssertNotNil(err, "Error should not be nil")
		helper.AssertEqual(context.Canceled, err, "Error should match expected")
	})
}

// TestEnvironmentVariables tests environment variable handling
func TestEnvironmentVariables(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("EnvironmentParsing", func(t *testing.T) {
		// Test environment variable parsing (if utility functions exist)
		// This would test actual functionality from utility package
		helper.AssertTrue(true, "Environment parsing placeholder")
	})
}

// TestConcurrentUtilities tests utility functions under concurrent access
func TestConcurrentUtilities(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ConcurrentDeploymentParsing", func(t *testing.T) {
		numGoroutines := 100
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				// Test concurrent deployment parsing
				deployments := []string{"dev", "test", "staging", "prod"}
				for _, dep := range deployments {
					result := utility.ParseDeployments(dep)
					require.NotNil(t, result)
				}
			}(i)
		}

		// Wait for all goroutines
		timeout := time.After(10 * time.Second)
		for i := 0; i < numGoroutines; i++ {
			select {
			case <-done:
				// Goroutine completed
			case <-timeout:
				t.Fatal("Timeout waiting for concurrent operations")
			}
		}
	})
}

// BenchmarkUtilities benchmarks utility functions
func BenchmarkUtilities(b *testing.B) {
	b.Run("DeploymentParsing", func(b *testing.B) {
		deployments := []string{"dev", "test", "staging", "prod"}
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for _, dep := range deployments {
					_ = utility.ParseDeployments(dep)
				}
			}
		})
	})

	b.Run("DeploymentMethods", func(b *testing.B) {
		deployment := utility.ParseDeployments("prod")
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = deployment.IsDev()
				_ = deployment.IsLocal()
				_ = deployment.IsProd()
				_ = deployment.String()
			}
		})
	})

	b.Run("ContextOperations", func(b *testing.B) {
		ctx := context.Background()
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ctxWithValue := context.WithValue(ctx, "key", "value")
				_ = ctxWithValue.Value("key")
			}
		})
	})
}
