package manapool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_UpdateSellerAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/account" {
			t.Fatalf("path = %s, want /account", r.URL.Path)
		}
		var payload SellerAccountUpdate
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload.SinglesLive == nil || !*payload.SinglesLive {
			t.Fatalf("singles_live not set")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"username":"test","email":"test@example.com","verified":true,"singles_live":true,"sealed_live":false,"payouts_enabled":true}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()
	value := true
	account, err := client.UpdateSellerAccount(ctx, SellerAccountUpdate{SinglesLive: &value})
	if err != nil {
		t.Fatalf("UpdateSellerAccount error: %v", err)
	}
	if !account.SinglesLive {
		t.Fatalf("singles_live = false, want true")
	}
}
