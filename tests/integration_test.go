package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCoreClient simula o Core Service para testes isolados
type MockCoreClient struct {
	mock.Mock
}

func (m *MockCoreClient) GetOrder(orderID string) (*CoreResponse, error) {
	args := m.Called(orderID)
	return args.Get(0).(*CoreResponse), args.Error(1)
}

func (m *MockCoreClient) UpdateOrderPaymentStatus(orderID, status string) error {
	args := m.Called(orderID, status)
	return args.Error(0)
}

// MockOperationClient simula o Operation Service para testes isolados
type MockOperationClient struct {
	mock.Mock
}

func (m *MockOperationClient) NotifyPaymentCompleted(orderID string, paymentData PaymentCompletedData) error {
	args := m.Called(orderID, paymentData)
	return args.Error(0)
}

// MockMercadoPagoClient simula o MercadoPago para testes isolados
type MockMercadoPagoClient struct {
	mock.Mock
}

func (m *MockMercadoPagoClient) CreatePayment(orderID string, amount float64) (*MercadoPagoResponse, error) {
	args := m.Called(orderID, amount)
	return args.Get(0).(*MercadoPagoResponse), args.Error(1)
}

func (m *MockMercadoPagoClient) GetPayment(paymentID string) (*MercadoPagoResponse, error) {
	args := m.Called(paymentID)
	return args.Get(0).(*MercadoPagoResponse), args.Error(1)
}

// CoreResponse representa resposta do Core Service
type CoreResponse struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
}

// OrderItem representa item do pedido
type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// PaymentCompletedData dados para notificar produção
type PaymentCompletedData struct {
	PaymentID string    `json:"payment_id"`
	Amount    float64   `json:"amount"`
	PaidAt    time.Time `json:"paid_at"`
}

// MercadoPagoResponse resposta do MercadoPago
type MercadoPagoResponse struct {
	ID                string    `json:"id"`
	Status            string    `json:"status"`
	ExternalReference string    `json:"external_reference"`
	QRCode            string    `json:"qr_code"`
	QRCodeBase64      string    `json:"qr_code_base64"`
	TransactionAmount float64   `json:"transaction_amount"`
	CreatedAt         time.Time `json:"date_created"`
	ExpiresAt         time.Time `json:"date_of_expiration"`
}

// TestPaymentCreationWithOrderValidation testa criação de pagamento com validação de pedido mockada
func TestPaymentCreationWithOrderValidation(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mocks
	mockCoreClient := new(MockCoreClient)
	mockMercadoPagoClient := new(MockMercadoPagoClient)

	// Mock responses
	expectedOrder := &CoreResponse{
		ID:          "order_123",
		CustomerID:  "customer_123",
		TotalAmount: 99.90,
		Status:      "pending",
		Items: []OrderItem{
			{ProductID: "prod_1", ProductName: "Hamburger", Quantity: 2, UnitPrice: 29.90, TotalPrice: 59.80},
			{ProductID: "prod_2", ProductName: "Fries", Quantity: 1, UnitPrice: 15.90, TotalPrice: 15.90},
		},
		CreatedAt: time.Now(),
	}

	expectedMPResponse := &MercadoPagoResponse{
		ID:                "mp_payment_123",
		Status:            "pending",
		ExternalReference: "order_123",
		QRCode:            "00020126580014BR.GOV.BCB.PIX...",
		QRCodeBase64:      "data:image/png;base64,iVBORw0KGgo...",
		TransactionAmount: 99.90,
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(30 * time.Minute),
	}

	// Configurar expectativas dos mocks
	mockCoreClient.On("GetOrder", "order_123").Return(expectedOrder, nil)
	mockMercadoPagoClient.On("CreatePayment", "order_123", 99.90).Return(expectedMPResponse, nil)

	// Rota de criação de pagamento
	router.POST("/payments", func(c *gin.Context) {
		var paymentRequest struct {
			OrderID string  `json:"order_id"`
			Amount  float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&paymentRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Validar pedido com Core Service (mockado)
		orderResponse, err := mockCoreClient.GetOrder(paymentRequest.OrderID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}

		// Validar valor
		if paymentRequest.Amount != orderResponse.TotalAmount {
			c.JSON(400, gin.H{"error": "Amount mismatch"})
			return
		}

		// Criar pagamento no MercadoPago (mockado)
		mpResponse, err := mockMercadoPagoClient.CreatePayment(paymentRequest.OrderID, paymentRequest.Amount)
		if err != nil {
			c.JSON(500, gin.H{"error": "Payment creation failed"})
			return
		}

		// Simular salvamento no MongoDB local
		paymentID := "payment_" + time.Now().Format("20060102150405")

		// Resposta
		c.JSON(200, gin.H{
			"payment_id":     paymentID,
			"order_id":       paymentRequest.OrderID,
			"amount":         paymentRequest.Amount,
			"status":         "pending",
			"qr_code":        mpResponse.QRCode,
			"qr_code_base64": mpResponse.QRCodeBase64,
			"mercadopago_id": mpResponse.ID,
			"expires_at":     mpResponse.ExpiresAt,
			"created_at":     time.Now(),
		})
	})

	// Teste
	paymentData := map[string]interface{}{
		"order_id": "order_123",
		"amount":   99.90,
	}

	jsonData, _ := json.Marshal(paymentData)
	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response["payment_id"])
	assert.Equal(t, "order_123", response["order_id"])
	assert.Equal(t, 99.90, response["amount"])
	assert.Equal(t, "pending", response["status"])
	assert.Equal(t, "mp_payment_123", response["mercadopago_id"])
	assert.NotEmpty(t, response["qr_code"])
	assert.NotEmpty(t, response["qr_code_base64"])

	// Verificar mocks
	mockCoreClient.AssertExpectations(t)
	mockMercadoPagoClient.AssertExpectations(t)
}

// TestPaymentWebhookProcessing testa processamento de webhook com notificações mockadas
func TestPaymentWebhookProcessing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mocks
	mockCoreClient := new(MockCoreClient)
	mockOperationClient := new(MockOperationClient)
	mockMercadoPagoClient := new(MockMercadoPagoClient)

	// Mock responses
	expectedMPResponse := &MercadoPagoResponse{
		ID:                "mp_payment_123",
		Status:            "approved",
		ExternalReference: "order_123",
		TransactionAmount: 99.90,
		CreatedAt:         time.Now(),
	}

	// Configurar expectativas - usar MatchedBy para timestamps flexíveis
	mockMercadoPagoClient.On("GetPayment", "mp_payment_123").Return(expectedMPResponse, nil)
	mockCoreClient.On("UpdateOrderPaymentStatus", "order_123", "paid").Return(nil)
	mockOperationClient.On("NotifyPaymentCompleted", "order_123", mock.MatchedBy(func(data PaymentCompletedData) bool {
		// Verificar se os dados estão corretos, ignorando diferenças mínimas de timestamp
		return data.PaymentID == "payment_123" &&
			data.Amount == 99.9 &&
			time.Since(data.PaidAt) < time.Second*2 // Aceita timestamps até 2 segundos de diferença
	})).Return(nil)

	// Rota de webhook
	router.POST("/webhook/payment/check", func(c *gin.Context) {
		var webhookData struct {
			Resource string `json:"resource"`
			Topic    string `json:"topic"`
			Type     string `json:"type"`
		}

		if err := c.ShouldBindJSON(&webhookData); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Extrair payment ID da URL do resource
		mpPaymentID := "mp_payment_123" // Simulado para teste

		// Consultar status no MercadoPago (mockado)
		mpResponse, err := mockMercadoPagoClient.GetPayment(mpPaymentID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get payment status"})
			return
		}

		if mpResponse.Status == "approved" {
			orderID := mpResponse.ExternalReference

			// Atualizar Core Service (mockado)
			err = mockCoreClient.UpdateOrderPaymentStatus(orderID, "paid")
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to update order"})
				return
			}

			// Notificar Operation Service (mockado)
			paymentData := PaymentCompletedData{
				PaymentID: "payment_123", // Simulado
				Amount:    mpResponse.TransactionAmount,
				PaidAt:    time.Now(),
			}

			err = mockOperationClient.NotifyPaymentCompleted(orderID, paymentData)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to notify operation"})
				return
			}

			c.JSON(200, gin.H{
				"status":         "processed",
				"order_id":       orderID,
				"payment_status": "approved",
			})
		} else {
			c.JSON(200, gin.H{
				"status":         "ignored",
				"payment_status": mpResponse.Status,
			})
		}
	})

	// Teste
	webhookData := map[string]interface{}{
		"resource": "https://api.mercadopago.com/v1/payments/mp_payment_123",
		"topic":    "payment",
		"type":     "payment",
	}

	jsonData, _ := json.Marshal(webhookData)
	req, _ := http.NewRequest("POST", "/webhook/payment/check", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "processed", response["status"])
	assert.Equal(t, "order_123", response["order_id"])
	assert.Equal(t, "approved", response["payment_status"])

	// Verificar mocks
	mockMercadoPagoClient.AssertExpectations(t)
	mockCoreClient.AssertExpectations(t)
	mockOperationClient.AssertExpectations(t)
}

// TestPaymentStatusCheck testa consulta de status sem dependências externas
func TestPaymentStatusCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simular dados locais do MongoDB
	payments := map[string]interface{}{
		"payment_123": map[string]interface{}{
			"payment_id":     "payment_123",
			"order_id":       "order_123",
			"amount":         99.90,
			"status":         "paid",
			"mercadopago_id": "mp_payment_123",
			"created_at":     time.Now().Add(-1 * time.Hour),
			"paid_at":        time.Now().Add(-30 * time.Minute),
		},
	}

	router.GET("/payments/:id", func(c *gin.Context) {
		paymentID := c.Param("id")

		payment, exists := payments[paymentID]
		if !exists {
			c.JSON(404, gin.H{"error": "Payment not found"})
			return
		}

		c.JSON(200, payment)
	})

	// Teste de pagamento existente
	req, _ := http.NewRequest("GET", "/payments/payment_123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "payment_123", response["payment_id"])
	assert.Equal(t, "order_123", response["order_id"])
	assert.Equal(t, 99.90, response["amount"])
	assert.Equal(t, "paid", response["status"])

	// Teste de pagamento não encontrado
	req, _ = http.NewRequest("GET", "/payments/nonexistent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

// TestPaymentListingWithFilters testa listagem com filtros sem dependências externas
func TestPaymentListingWithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock de dados locais
	payments := []map[string]interface{}{
		{
			"payment_id": "payment_123",
			"order_id":   "order_123",
			"status":     "paid",
			"amount":     99.90,
			"created_at": time.Now().Add(-2 * time.Hour),
		},
		{
			"payment_id": "payment_456",
			"order_id":   "order_456",
			"status":     "pending",
			"amount":     45.50,
			"created_at": time.Now().Add(-1 * time.Hour),
		},
		{
			"payment_id": "payment_789",
			"order_id":   "order_789",
			"status":     "failed",
			"amount":     75.00,
			"created_at": time.Now().Add(-30 * time.Minute),
		},
	}

	router.GET("/payments", func(c *gin.Context) {
		status := c.Query("status")

		var filteredPayments []map[string]interface{}

		for _, payment := range payments {
			if status == "" || payment["status"] == status {
				filteredPayments = append(filteredPayments, payment)
			}
		}

		c.JSON(200, gin.H{
			"payments": filteredPayments,
			"total":    len(filteredPayments),
		})
	})

	// Teste sem filtro
	req, _ := http.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	paymentsResponse := response["payments"].([]interface{})
	assert.Len(t, paymentsResponse, 3)
	assert.Equal(t, float64(3), response["total"])

	// Teste com filtro de status
	req, _ = http.NewRequest("GET", "/payments?status=paid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	paymentsResponse = response["payments"].([]interface{})
	assert.Len(t, paymentsResponse, 1)
	assert.Equal(t, float64(1), response["total"])

	firstPayment := paymentsResponse[0].(map[string]interface{})
	assert.Equal(t, "paid", firstPayment["status"])
}
