package gateways

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type mockResponse struct {
	Message string `json:"message"`
}

func TestMercadoPagoClientRest_Get(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expected     string
		expectError  bool
	}{
		{
			name:         "success",
			statusCode:   http.StatusOK,
			responseBody: `{"message":"hello get"}`,
			expected:     "hello get",
			expectError:  false,
		},
		{
			name:         "invalid json",
			statusCode:   http.StatusOK,
			responseBody: `{invalid`,
			expected:     "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := &MercadoPagoClientRest{
				client: resty.New().SetBaseURL(server.URL),
			}

			var resp mockResponse
			_, err := client.Get("/", &resp)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp.Message)
			}
		})
	}
}

func TestMercadoPagoClientRest_Post(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expected     string
		expectError  bool
	}{
		{
			name:         "success",
			statusCode:   http.StatusOK,
			responseBody: `{"message":"hello post"}`,
			expected:     "hello post",
			expectError:  false,
		},
		{
			name:         "invalid json",
			statusCode:   http.StatusOK,
			responseBody: `{not_json`,
			expected:     "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// opcional: verificar o corpo da requisição recebida
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := &MercadoPagoClientRest{
				client: resty.New().SetBaseURL(server.URL),
			}

			var result mockResponse
			_, err := client.Post("/", map[string]string{"input": "value"}, &result)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.Message)
			}
		})
	}
}
