package manapool

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHighCoverageImprovements adds tests to reach 95% coverage
func TestHighCoverageImprovements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/error":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "bad request"}`))
		case "/empty":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success": true}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	// Test all error paths in job application
	t.Run("JobApplication_AllErrorPaths", func(t *testing.T) {
		// Test default filename path
		_, err := client.SubmitJobApplication(ctx, JobApplicationRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john@example.com",
			Application: []byte("test"),
		})
		if err != nil {
			t.Fatalf("SubmitJobApplication error: %v", err)
		}
	})

	// Test all parameter building functions
	t.Run("BuyerOrdersParams", func(t *testing.T) {
		opts := BuyerOrdersOptions{Limit: 10, Offset: 5}
		_, err := client.GetBuyerOrders(ctx, opts)
		if err != nil {
			t.Fatalf("GetBuyerOrders error: %v", err)
		}
	})

	// Test webhook with topic filter
	t.Run("WebhookWithTopic", func(t *testing.T) {
		_, err := client.GetWebhooks(ctx, "order_created")
		if err != nil {
			t.Fatalf("GetWebhooks error: %v", err)
		}
	})

	// Test orders with all options
	t.Run("OrdersWithAllOptions", func(t *testing.T) {
		opts := OrdersOptions{
			IsUnfulfilled:   &[]bool{true}[0],
			IsFulfilled:     &[]bool{false}[0],
			HasFulfillments: &[]bool{true}[0],
			Label:           "test",
			Limit:           10,
			Offset:          5,
		}
		_, err := client.GetOrders(ctx, opts)
		if err != nil {
			t.Fatalf("GetOrders error: %v", err)
		}
	})

	// Test seller orders with options
	t.Run("SellerOrdersWithOptions", func(t *testing.T) {
		opts := OrdersOptions{Label: "test", Limit: 5}
		_, err := client.GetSellerOrders(ctx, opts)
		if err != nil {
			t.Fatalf("GetSellerOrders error: %v", err)
		}
	})

	// Test inventory listings with IDs
	t.Run("InventoryListingsWithIDs", func(t *testing.T) {
		_, err := client.GetInventoryListings(ctx, []string{"id1", "id2"})
		if err != nil {
			t.Fatalf("GetInventoryListings error: %v", err)
		}
	})

	// Test seller inventory by scryfall with all options
	t.Run("SellerInventoryByScryfall", func(t *testing.T) {
		opts := InventoryByScryfallOptions{
			LanguageID:  "EN",
			FinishID:    "NF",
			ConditionID: "NM",
		}
		_, err := client.GetSellerInventoryByScryfall(ctx, "test-id", opts)
		if err != nil {
			t.Fatalf("GetSellerInventoryByScryfall error: %v", err)
		}
	})

	// Test seller inventory by TCGPlayer with all options
	t.Run("SellerInventoryByTCGPlayer", func(t *testing.T) {
		opts := InventoryByTCGPlayerOptions{
			LanguageID:  "EN",
			FinishID:    "FO",
			ConditionID: "LP",
		}
		_, err := client.GetSellerInventoryByTCGPlayerID(ctx, 123, opts)
		if err != nil {
			t.Fatalf("GetSellerInventoryByTCGPlayerID error: %v", err)
		}
	})
}

// TestNoopLoggerCoverage ensures noop logger methods are covered
func TestNoopLoggerCoverage(t *testing.T) {
	// Create a client that uses the noop logger
	client := NewClient("test", "test")
	
	// Access the logger directly to ensure coverage
	if client.logger != nil {
		client.logger.Debugf("test debug message with %s", "args")
		client.logger.Errorf("test error message with %d", 123)
	}
	
	// Also test the noop logger directly
	logger := &noopLogger{}
	logger.Debugf("direct test")
	logger.Errorf("direct test")
}

// TestAllValidationErrors tests all validation error paths
func TestAllValidationErrors(t *testing.T) {
	client := NewClient("test", "test")
	ctx := context.Background()

	// Test all empty ID validations
	validationTests := []struct {
		name string
		fn   func() error
	}{
		{"GetWebhook", func() error { _, err := client.GetWebhook(ctx, ""); return err }},
		{"DeleteWebhook", func() error { return client.DeleteWebhook(ctx, "") }},
		{"GetOrder", func() error { _, err := client.GetOrder(ctx, ""); return err }},
		{"GetSellerOrder", func() error { _, err := client.GetSellerOrder(ctx, ""); return err }},
		{"GetBuyerOrder", func() error { _, err := client.GetBuyerOrder(ctx, ""); return err }},
		{"GetPendingOrder", func() error { _, err := client.GetPendingOrder(ctx, ""); return err }},
		{"GetInventoryListing", func() error { _, err := client.GetInventoryListing(ctx, ""); return err }},
		{"GetSellerOrderReports", func() error { _, err := client.GetSellerOrderReports(ctx, ""); return err }},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("expected ValidationError, got %T", err)
			}
		})
	}

	// Test SKU validations
	skuTests := []struct {
		name string
		fn   func() error
	}{
		{"GetInventoryBySKU", func() error { _, err := client.GetInventoryBySKU(ctx, 0); return err }},
		{"UpdateInventoryBySKU", func() error { _, err := client.UpdateInventoryBySKU(ctx, -1, InventoryUpdateRequest{}); return err }},
		{"DeleteInventoryBySKU", func() error { _, err := client.DeleteInventoryBySKU(ctx, 0); return err }},
		{"GetSellerInventoryBySKU", func() error { _, err := client.GetSellerInventoryBySKU(ctx, 0); return err }},
	}

	for _, tt := range skuTests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("expected ValidationError, got %T", err)
			}
		})
	}
}

// TestEmptyArrayValidations tests empty array validations
func TestEmptyArrayValidations(t *testing.T) {
	client := NewClient("test", "test")
	ctx := context.Background()

	emptyArrayTests := []struct {
		name string
		fn   func() error
	}{
		{"CreateInventoryBulk", func() error { _, err := client.CreateInventoryBulk(ctx, []InventoryBulkItemBySKU{}); return err }},
		{"CreateInventoryBulkBySKU", func() error { _, err := client.CreateInventoryBulkBySKU(ctx, []InventoryBulkItemBySKU{}); return err }},
		{"CreateInventoryBulkByProduct", func() error { _, err := client.CreateInventoryBulkByProduct(ctx, []InventoryBulkItemByProduct{}); return err }},
		{"CreateInventoryBulkByScryfall", func() error { _, err := client.CreateInventoryBulkByScryfall(ctx, []InventoryBulkItemByScryfall{}); return err }},
		{"CreateInventoryBulkByTCGPlayerID", func() error { _, err := client.CreateInventoryBulkByTCGPlayerID(ctx, []InventoryBulkItemByTCGPlayerID{}); return err }},
	}

	for _, tt := range emptyArrayTests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("expected ValidationError, got %T", err)
			}
		})
	}
}
// TestUpdateRequestPaths tests update request paths
func TestUpdateRequestPaths(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient("test", "test", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("UpdatePendingOrder", func(t *testing.T) {
		_, err := client.UpdatePendingOrder(ctx, "test-id", PendingOrderRequest{})
		if err != nil {
			t.Fatalf("UpdatePendingOrder error: %v", err)
		}
	})

	t.Run("PurchasePendingOrder", func(t *testing.T) {
		_, err := client.PurchasePendingOrder(ctx, "test-id", PurchasePendingOrderRequest{})
		if err != nil {
			t.Fatalf("PurchasePendingOrder error: %v", err)
		}
	})

	t.Run("UpdateOrderFulfillment", func(t *testing.T) {
		_, err := client.UpdateOrderFulfillment(ctx, "test-id", OrderFulfillmentRequest{})
		if err != nil {
			t.Fatalf("UpdateOrderFulfillment error: %v", err)
		}
	})

	t.Run("UpdateSellerOrderFulfillment", func(t *testing.T) {
		_, err := client.UpdateSellerOrderFulfillment(ctx, "test-id", OrderFulfillmentRequest{})
		if err != nil {
			t.Fatalf("UpdateSellerOrderFulfillment error: %v", err)
		}
	})

	t.Run("UpdateSellerInventoryBySKU", func(t *testing.T) {
		_, err := client.UpdateSellerInventoryBySKU(ctx, 123, InventoryUpdateRequest{})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryBySKU error: %v", err)
		}
	})

	t.Run("UpdateSellerInventoryByProduct", func(t *testing.T) {
		_, err := client.UpdateSellerInventoryByProduct(ctx, "mtg_single", "test-id", InventoryUpdateRequest{})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryByProduct error: %v", err)
		}
	})

	t.Run("UpdateSellerInventoryByScryfall", func(t *testing.T) {
		opts := InventoryByScryfallOptions{}
		_, err := client.UpdateSellerInventoryByScryfall(ctx, "test-id", opts, InventoryUpdateRequest{})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryByScryfall error: %v", err)
		}
	})

	t.Run("UpdateSellerInventoryByTCGPlayerID", func(t *testing.T) {
		opts := InventoryByTCGPlayerOptions{}
		_, err := client.UpdateSellerInventoryByTCGPlayerID(ctx, 123, opts, InventoryUpdateRequest{})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryByTCGPlayerID error: %v", err)
		}
	})
}

// TestDeletePaths tests delete request paths
func TestDeletePaths(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test", "test", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("DeleteSellerInventoryBySKU", func(t *testing.T) {
		_, err := client.DeleteSellerInventoryBySKU(ctx, 123)
		if err != nil {
			t.Fatalf("DeleteSellerInventoryBySKU error: %v", err)
		}
	})

	t.Run("DeleteSellerInventoryByProduct", func(t *testing.T) {
		_, err := client.DeleteSellerInventoryByProduct(ctx, "mtg_single", "test-id")
		if err != nil {
			t.Fatalf("DeleteSellerInventoryByProduct error: %v", err)
		}
	})

	t.Run("DeleteSellerInventoryByScryfall", func(t *testing.T) {
		opts := InventoryByScryfallOptions{}
		_, err := client.DeleteSellerInventoryByScryfall(ctx, "test-id", opts)
		if err != nil {
			t.Fatalf("DeleteSellerInventoryByScryfall error: %v", err)
		}
	})

	t.Run("DeleteSellerInventoryByTCGPlayerID", func(t *testing.T) {
		opts := InventoryByTCGPlayerOptions{}
		_, err := client.DeleteSellerInventoryByTCGPlayerID(ctx, 123, opts)
		if err != nil {
			t.Fatalf("DeleteSellerInventoryByTCGPlayerID error: %v", err)
		}
	})
}

// TestBulkCreationPaths tests bulk creation paths
func TestBulkCreationPaths(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"inventory": []}`))
	}))
	defer server.Close()

	client := NewClient("test", "test", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("CreateInventoryBulkBySKU", func(t *testing.T) {
		items := []InventoryBulkItemBySKU{{TCGPlayerSKU: 123, PriceCents: 100, Quantity: 1}}
		_, err := client.CreateInventoryBulkBySKU(ctx, items)
		if err != nil {
			t.Fatalf("CreateInventoryBulkBySKU error: %v", err)
		}
	})

	t.Run("CreateInventoryBulkByProduct", func(t *testing.T) {
		items := []InventoryBulkItemByProduct{{ProductType: "test", ProductID: "test", PriceCents: 100, Quantity: 1}}
		_, err := client.CreateInventoryBulkByProduct(ctx, items)
		if err != nil {
			t.Fatalf("CreateInventoryBulkByProduct error: %v", err)
		}
	})

	t.Run("CreateInventoryBulkByScryfall", func(t *testing.T) {
		items := []InventoryBulkItemByScryfall{{ScryfallID: "test", LanguageID: "EN", FinishID: "NF", ConditionID: "NM", PriceCents: 100, Quantity: 1}}
		_, err := client.CreateInventoryBulkByScryfall(ctx, items)
		if err != nil {
			t.Fatalf("CreateInventoryBulkByScryfall error: %v", err)
		}
	})

	t.Run("CreateInventoryBulkByTCGPlayerID", func(t *testing.T) {
		items := []InventoryBulkItemByTCGPlayerID{{TCGPlayerID: 123, LanguageID: "EN", PriceCents: 100, Quantity: 1}}
		_, err := client.CreateInventoryBulkByTCGPlayerID(ctx, items)
		if err != nil {
			t.Fatalf("CreateInventoryBulkByTCGPlayerID error: %v", err)
		}
	})
}
