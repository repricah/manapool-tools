package manapool

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// OptimizeCart creates an optimized cart.
func (c *Client) OptimizeCart(ctx context.Context, req OptimizerRequest) (*OptimizedCart, error) {
	resp, err := c.doJSONRequest(ctx, "POST", "/buyer/optimizer", nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to optimize cart: %w", err)
	}

	var cart OptimizedCart
	if err := c.decodeResponse(resp, &cart); err != nil {
		return nil, fmt.Errorf("failed to decode optimized cart: %w", err)
	}

	return &cart, nil
}

// GetBuyerOrders retrieves buyer orders with optional filtering.
func (c *Client) GetBuyerOrders(ctx context.Context, opts BuyerOrdersOptions) (*BuyerOrdersResponse, error) {
	params := url.Values{}
	if opts.Since != nil {
		params.Add("since", opts.Since.Format(time.RFC3339Nano))
	}
	if opts.Limit > 0 {
		params.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Add("offset", strconv.Itoa(opts.Offset))
	}

	resp, err := c.doRequest(ctx, "GET", "/buyer/orders", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get buyer orders: %w", err)
	}

	var orders BuyerOrdersResponse
	if err := c.decodeResponse(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode buyer orders: %w", err)
	}

	return &orders, nil
}

// GetBuyerOrder retrieves a buyer order by ID.
func (c *Client) GetBuyerOrder(ctx context.Context, id string) (*BuyerOrderResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/buyer/orders/%s", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get buyer order: %w", err)
	}

	var order BuyerOrderResponse
	if err := c.decodeResponse(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to decode buyer order: %w", err)
	}

	return &order, nil
}

// CreatePendingOrder creates a pending order.
func (c *Client) CreatePendingOrder(ctx context.Context, req PendingOrderRequest) (*PendingOrder, error) {
	resp, err := c.doJSONRequest(ctx, "POST", "/buyer/orders/pending-orders", nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending order: %w", err)
	}

	var pending PendingOrder
	if err := c.decodeResponse(resp, &pending); err != nil {
		return nil, fmt.Errorf("failed to decode pending order: %w", err)
	}

	return &pending, nil
}

// GetPendingOrder retrieves a pending order by ID.
func (c *Client) GetPendingOrder(ctx context.Context, id string) (*PendingOrder, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/buyer/orders/pending-orders/%s", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending order: %w", err)
	}

	var pending PendingOrder
	if err := c.decodeResponse(resp, &pending); err != nil {
		return nil, fmt.Errorf("failed to decode pending order: %w", err)
	}

	return &pending, nil
}

// UpdatePendingOrder updates a pending order.
func (c *Client) UpdatePendingOrder(ctx context.Context, id string, req PendingOrderRequest) (*PendingOrder, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/buyer/orders/pending-orders/%s", id)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update pending order: %w", err)
	}

	var pending PendingOrder
	if err := c.decodeResponse(resp, &pending); err != nil {
		return nil, fmt.Errorf("failed to decode pending order: %w", err)
	}

	return &pending, nil
}

// PurchasePendingOrder purchases a pending order.
func (c *Client) PurchasePendingOrder(ctx context.Context, id string, req PurchasePendingOrderRequest) (*PendingOrder, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/buyer/orders/pending-orders/%s/purchase", id)
	resp, err := c.doJSONRequest(ctx, "POST", endpoint, nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to purchase pending order: %w", err)
	}

	var pending PendingOrder
	if err := c.decodeResponse(resp, &pending); err != nil {
		return nil, fmt.Errorf("failed to decode purchased order: %w", err)
	}

	return &pending, nil
}

// GetBuyerCredit retrieves buyer credit balance.
func (c *Client) GetBuyerCredit(ctx context.Context) (*BuyerCredit, error) {
	resp, err := c.doRequest(ctx, "GET", "/buyer/credit", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get buyer credit: %w", err)
	}

	var credit BuyerCredit
	if err := c.decodeResponse(resp, &credit); err != nil {
		return nil, fmt.Errorf("failed to decode buyer credit: %w", err)
	}

	return &credit, nil
}
