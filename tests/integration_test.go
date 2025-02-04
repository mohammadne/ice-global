package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	apiBaseURL = "http://localhost:8080" // Adjust based on your server config
	client     = &http.Client{Timeout: 5 * time.Second}
	cookie     = &http.Cookie{
		Name:  "ice_session_id",
		Value: "test-session-id",
	}
)

// waitForService waits until the service is available
func waitForService(*testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/healthz/liveness", apiBaseURL), nil)
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// Scenario 1: Health Check - Liveness
func TestHealthCheck_Liveness(t *testing.T) {
	waitForService(t)

	resp, err := client.Get(fmt.Sprintf("%s/healthz/liveness", apiBaseURL))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Scenario 2: Health Check - Readiness
func TestHealthCheck_Readiness(t *testing.T) {
	waitForService(t)

	resp, err := client.Get(fmt.Sprintf("%s/healthz/readiness", apiBaseURL))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Scenario 3: Add Item to Cart (POST /add-item)
func TestAddItemToCart(t *testing.T) {
	// Assuming cart exists for the session
	item := struct {
		ItemId   int    `form:"item_id" binding:"required"`
		Quantity string `form:"quantity" binding:"required"`
	}{
		ItemId:   1,
		Quantity: "2",
	}

	payload, _ := json.Marshal(item)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/add-item", apiBaseURL), bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode) // Expecting a redirect
}

// Scenario 4: Show Add Item Form (GET /)
func TestShowAddItemForm(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/", apiBaseURL), nil)
	require.NoError(t, err)
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Items") // Check if items are present in the response
}

// Scenario 5: Remove Item from Cart (GET /remove-cart-item)
func TestRemoveItemFromCart(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/remove-cart-item?cart_item_id=1", apiBaseURL), nil)
	require.NoError(t, err)
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode) // Expecting a redirect
}
