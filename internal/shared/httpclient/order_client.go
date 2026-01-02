package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type OrderClient struct {
	baseURL string
	client  *http.Client
}

type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// addServiceAuth adds service authentication headers to HTTP requests
func (c *OrderClient) addServiceAuth(req *http.Request) {
	// Add service-to-service authentication headers
	req.Header.Set("X-Service-Name", "payment-service")
	if apiKey := os.Getenv("PAYMENT_SERVICE_API_KEY"); apiKey != "" {
		req.Header.Set("X-Service-Key", apiKey)
	}
}

func (c *OrderClient) FindByID(ctx context.Context, orderID string) (Order, error) {
	url := fmt.Sprintf("%s/admin/orders/%s", c.baseURL, orderID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Order{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Add service authentication
	c.addServiceAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return Order{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Order{}, fmt.Errorf("failed to get order: status %d", resp.StatusCode)
	}

	var order Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return Order{}, err
	}

	return order, nil
}

func (c *OrderClient) Update(ctx context.Context, order Order) (Order, error) {
	url := fmt.Sprintf("%s/admin/orders/%s", c.baseURL, order.ID)

	payload := map[string]string{
		"status": order.Status,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Order{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Order{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Add service authentication
	c.addServiceAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return Order{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Order{}, fmt.Errorf("failed to update order: status %d", resp.StatusCode)
	}

	var updatedOrder Order
	if err := json.NewDecoder(resp.Body).Decode(&updatedOrder); err != nil {
		return Order{}, err
	}

	return updatedOrder, nil
}
