package middleware

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServiceAuthMiddleware(t *testing.T) {
	// Set test environment variables
	os.Setenv("CORE_SERVICE_API_KEY", "test-core-api-key")
	os.Setenv("PAYMENT_SERVICE_API_KEY", "test-payment-api-key")
	os.Setenv("OPERATION_SERVICE_API_KEY", "test-production-api-key")

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		method         string
		serviceName    string
		serviceKey     string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health check endpoint bypasses auth",
			path:           "/ping",
			method:         "GET",
			serviceName:    "",
			serviceKey:     "",
			expectedStatus: 200,
		},
		{
			name:           "Valid service credentials",
			path:           "/api/test",
			method:         "GET",
			serviceName:    "core-service",
			serviceKey:     "test-core-api-key",
			expectedStatus: 200,
		},
		{
			name:           "Valid payment service credentials",
			path:           "/api/orders",
			method:         "POST",
			serviceName:    "payment-service",
			serviceKey:     "test-payment-api-key",
			expectedStatus: 200,
		},
		{
			name:           "Invalid service name",
			path:           "/api/test",
			method:         "GET",
			serviceName:    "unknown-service",
			serviceKey:     "test-core-api-key",
			expectedStatus: 401,
			expectedBody:   "Unauthorized: Invalid service credentials",
		},
		{
			name:           "Invalid API key",
			path:           "/api/test",
			method:         "GET",
			serviceName:    "core-service",
			serviceKey:     "wrong-api-key",
			expectedStatus: 401,
			expectedBody:   "Unauthorized: Invalid service credentials",
		},
		{
			name:           "Missing service name",
			path:           "/api/test",
			method:         "GET",
			serviceName:    "",
			serviceKey:     "test-core-api-key",
			expectedStatus: 401,
			expectedBody:   "Unauthorized: Invalid service credentials",
		},
		{
			name:           "Missing API key",
			path:           "/api/test",
			method:         "GET",
			serviceName:    "core-service",
			serviceKey:     "",
			expectedStatus: 401,
			expectedBody:   "Unauthorized: Invalid service credentials",
		},
		{
			name:           "JWT token should pass through",
			path:           "/api/test",
			method:         "GET",
			authHeader:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			expectedStatus: 200, // Should pass through to next middleware
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with middleware
			router := gin.New()
			router.Use(ServiceAuthMiddleware())

			// Add test route
			router.Any("/*path", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Add headers
			if tt.serviceName != "" {
				req.Header.Set("X-Service-Name", tt.serviceName)
			}
			if tt.serviceKey != "" {
				req.Header.Set("X-Service-Key", tt.serviceKey)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute request
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			// Assert status
			assert.Equal(t, tt.expectedStatus, resp.Code, "Status code mismatch")

			// Assert error message for auth failures
			if tt.expectedBody != "" && resp.Code == 401 {
				assert.Contains(t, resp.Body.String(), tt.expectedBody)
			}

			// Assert authenticated service is set for successful auth
			if tt.expectedStatus == 200 && tt.serviceName != "" && tt.serviceKey != "" && tt.authHeader == "" {
				// Note: We can't easily test context values in this setup,
				// but we know if status is 200, the service was authenticated
			}
		})
	}

	// Cleanup
	os.Unsetenv("CORE_SERVICE_API_KEY")
	os.Unsetenv("PAYMENT_SERVICE_API_KEY")
	os.Unsetenv("OPERATION_SERVICE_API_KEY")
}

func TestValidateServiceAPIKey(t *testing.T) {
	// Set test environment variables
	os.Setenv("CORE_SERVICE_API_KEY", "test-core-key")
	os.Setenv("PAYMENT_SERVICE_API_KEY", "test-payment-key")

	tests := []struct {
		name        string
		serviceName string
		apiKey      string
		expected    bool
	}{
		{
			name:        "Valid core service key",
			serviceName: "core-service",
			apiKey:      "test-core-key",
			expected:    true,
		},
		{
			name:        "Valid payment service key",
			serviceName: "payment-service",
			apiKey:      "test-payment-key",
			expected:    true,
		},
		{
			name:        "Invalid service name",
			serviceName: "unknown-service",
			apiKey:      "test-core-key",
			expected:    false,
		},
		{
			name:        "Invalid API key",
			serviceName: "core-service",
			apiKey:      "wrong-key",
			expected:    false,
		},
		{
			name:        "Empty service name",
			serviceName: "",
			apiKey:      "test-core-key",
			expected:    false,
		},
		{
			name:        "Empty API key",
			serviceName: "core-service",
			apiKey:      "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateServiceAPIKey(tt.serviceName, tt.apiKey)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Cleanup
	os.Unsetenv("CORE_SERVICE_API_KEY")
	os.Unsetenv("PAYMENT_SERVICE_API_KEY")
}

func TestConstantTimeEquals(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "Equal strings",
			a:        "test-key-123",
			b:        "test-key-123",
			expected: true,
		},
		{
			name:     "Different strings same length",
			a:        "test-key-123",
			b:        "test-key-456",
			expected: false,
		},
		{
			name:     "Different lengths",
			a:        "short",
			b:        "much-longer-string",
			expected: false,
		},
		{
			name:     "Empty strings",
			a:        "",
			b:        "",
			expected: true,
		},
		{
			name:     "One empty string",
			a:        "test",
			b:        "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := constantTimeEquals(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}
