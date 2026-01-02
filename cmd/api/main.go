package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// @title           GoLunch Payment Service API
// @version         1.0
// @description     API para processamento de pagamentos da lanchonete GoLunch
// @host            localhost:8082
// @BasePath        /
func main() {
	r := gin.Default()

	// Default Routes
	r.GET("/ping", ping)

	// Payment Routes (simplificado para funcionar)
	r.POST("/payments", createPayment)
	r.POST("/webhook/payment/check", checkPayment)

	log.Println("Payment Service starting on port 8082...")
	r.Run(":8082")
}

// Ping godoc
// @Summary      Answers with "pong"
// @Description  Health Check
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200 {object}  map[string]string
// @Router       /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Create Payment endpoint (simples para testes)
func createPayment(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Simular criação de pagamento
	c.JSON(200, gin.H{
		"payment_id": "payment_123",
		"qr_code":    "mock-qr-code-data",
		"status":     "pending",
		"order_id":   request["order_id"],
	})
}

// Check Payment webhook (simples para testes)
func checkPayment(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Simular verificação de pagamento
	c.JSON(200, gin.H{
		"status":  "approved",
		"message": "Payment verified successfully",
	})
}

