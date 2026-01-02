package gateways

import (
	"context"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fiap-161/tc-golunch-payment-service/internal/shared"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func newMockedMercadoPagoClient(baseURL string) *MercadoPagoClient {
	client := &MercadoPagoClientRest{
		resty.New().
			SetBaseURL(baseURL).
			SetHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer dummy",
			}),
	}
	return &MercadoPagoClient{client: client}
}

func mockEnvVars() {
	viper.Set(shared.MercadoPagoQRCodePath, "/mocked-path/{user_id}/{external_pos_id}")
	os.Setenv("MERCADO_PAGO_SELLER_APP_USER_ID", "userid_12")
	os.Setenv("MERCADO_PAGO_EXTERNAL_POS_ID", "posid_34")
}

func startTestServer(t *testing.T, expectedPath string, statusCode int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedPath, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))
}

func TestGenerateQRCode(t *testing.T) {
	mockEnvVars()

	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expectError  bool
		expectedQR   string
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			responseBody: `{
				"in_store_order_id": "orderId12",
				"qr_data": "http://mocked-qr-code-data"
			}`,
			expectError: false,
			expectedQR:  "http://mocked-qr-code-data",
		},
		{
			name:         "error status code",
			statusCode:   http.StatusInternalServerError,
			responseBody: "",
			expectError:  true,
			expectedQR:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := startTestServer(t, "/mocked-path/userid_12/posid_34", tt.statusCode, tt.responseBody)
			defer server.Close()

			mpClient := newMockedMercadoPagoClient(server.URL)

			params := entities.GenerateQRCodeParams{
				OrderID: "orderId",
				Items: []entities.Item{
					{
						ID:          "itemId",
						Name:        "itemName",
						Price:       12.3,
						Description: "itemDescription",
						Quantity:    1,
						Amount:      12.3,
					},
				},
			}

			qrData, err := mpClient.GenerateQRCode(context.Background(), params)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, qrData)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedQR, qrData)
			}
		})
	}
}

func TestCheckPayment(t *testing.T) {
	tests := []struct {
		name                      string
		statusCode                int
		responseBody              string
		expectError               bool
		expectedOrderStatus       string
		expectedExternalReference string
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			responseBody: `{
				"external_reference": "orderId12",
				"order_status": "approved"
			}`,
			expectError:               false,
			expectedOrderStatus:       "approved",
			expectedExternalReference: "orderId12",
		},
		{
			name:                      "invalid JSON",
			statusCode:                http.StatusOK,
			responseBody:              `{invalid-json}`,
			expectError:               true,
			expectedOrderStatus:       "",
			expectedExternalReference: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := startTestServer(t, "/expected-path", tt.statusCode, tt.responseBody)
			defer server.Close()

			mpClient := newMockedMercadoPagoClient(server.URL)

			resp, err := mpClient.CheckPayment(context.Background(), "/expected-path")

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, resp.ExternalReference)
				assert.Empty(t, resp.OrderStatus)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExternalReference, resp.ExternalReference)
				assert.Equal(t, tt.expectedOrderStatus, resp.OrderStatus)
			}
		})
	}
}
