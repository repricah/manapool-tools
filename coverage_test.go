package manapool

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestCoverageImprovements contains additional tests to improve coverage
func TestCoverageImprovements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("BuyerOrdersOptions_toParams", func(t *testing.T) {
		since := Timestamp{time.Now()}
		opts := BuyerOrdersOptions{
			Since:  &since,
			Limit:  10,
			Offset: 5,
		}
		
		// Test the buildBuyerOrdersParams function indirectly
		_, err := client.GetBuyerOrders(ctx, opts)
		if err != nil {
			t.Fatalf("GetBuyerOrders error: %v", err)
		}
	})

	t.Run("OrdersOptions_buildParams", func(t *testing.T) {
		since := Timestamp{time.Now()}
		opts := OrdersOptions{
			Since:           &since,
			IsUnfulfilled:   boolPtr(true),
			IsFulfilled:     boolPtr(false),
			HasFulfillments: boolPtr(true),
			Label:           "test-label",
			Limit:           10,
			Offset:          5,
		}
		
		// Test the buildOrdersParams function indirectly
		_, err := client.GetOrders(ctx, opts)
		if err != nil {
			t.Fatalf("GetOrders error: %v", err)
		}
	})

	t.Run("InventoryByScryfallOptions_toParams", func(t *testing.T) {
		opts := InventoryByScryfallOptions{
			LanguageID:  "EN",
			FinishID:    "NF",
			ConditionID: "NM",
		}
		
		params := opts.toParams()
		if params.Get("language_id") != "EN" {
			t.Errorf("language_id = %s, want EN", params.Get("language_id"))
		}
		if params.Get("finish_id") != "NF" {
			t.Errorf("finish_id = %s, want NF", params.Get("finish_id"))
		}
		if params.Get("condition_id") != "NM" {
			t.Errorf("condition_id = %s, want NM", params.Get("condition_id"))
		}
	})

	t.Run("InventoryByTCGPlayerOptions_toParams", func(t *testing.T) {
		opts := InventoryByTCGPlayerOptions{
			LanguageID:  "EN",
			FinishID:    "FO",
			ConditionID: "LP",
		}
		
		params := opts.toParams()
		if params.Get("language_id") != "EN" {
			t.Errorf("language_id = %s, want EN", params.Get("language_id"))
		}
		if params.Get("finish_id") != "FO" {
			t.Errorf("finish_id = %s, want FO", params.Get("finish_id"))
		}
		if params.Get("condition_id") != "LP" {
			t.Errorf("condition_id = %s, want LP", params.Get("condition_id"))
		}
	})

	t.Run("EmptyOptions_toParams", func(t *testing.T) {
		opts := InventoryByScryfallOptions{}
		params := opts.toParams()
		
		// Should return empty params
		if len(params) != 0 {
			t.Errorf("expected empty params, got %v", params)
		}
	})
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}

// TestAdditionalErrorCases tests additional error scenarios
func TestAdditionalErrorCases(t *testing.T) {
	// Test with server that returns different status codes
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

	t.Run("doRequest_with_query_params", func(t *testing.T) {
		params := url.Values{}
		params.Add("test", "value")
		params.Add("another", "param")
		
		resp, err := client.doRequest(ctx, "GET", "/test", params)
		if err != nil {
			t.Fatalf("doRequest error: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code = %d, want 200", resp.StatusCode)
		}
	})

	t.Run("decodeResponse_with_empty_body", func(t *testing.T) {
		resp, err := client.doRequest(ctx, "GET", "/empty", nil)
		if err != nil {
			t.Fatalf("doRequest error: %v", err)
		}
		defer resp.Body.Close()
		
		// Should handle empty response body gracefully
		err = client.decodeResponse(resp, nil)
		if err != nil {
			t.Fatalf("decodeResponse error: %v", err)
		}
	})
}
