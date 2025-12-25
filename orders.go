package manapool

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GetOrders retrieves order summaries.
func (c *Client) GetOrders(ctx context.Context, opts OrdersOptions) (*OrdersResponse, error) {
	params := buildOrdersParams(opts)
	resp, err := c.doRequest(ctx, "GET", "/orders", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var orders OrdersResponse
	if err := c.decodeResponse(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %w", err)
	}

	return &orders, nil
}

// GetOrder retrieves order details by ID.
func (c *Client) GetOrder(ctx context.Context, id string) (*OrderDetailsResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/orders/%s", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var order OrderDetailsResponse
	if err := c.decodeResponse(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to decode order: %w", err)
	}

	return &order, nil
}

// UpdateOrderFulfillment updates the fulfillment for an order.
func (c *Client) UpdateOrderFulfillment(ctx context.Context, id string, req OrderFulfillmentRequest) (*OrderFulfillmentResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/orders/%s/fulfillment", id)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update order fulfillment: %w", err)
	}

	var fulfillment OrderFulfillmentResponse
	if err := c.decodeResponse(resp, &fulfillment); err != nil {
		return nil, fmt.Errorf("failed to decode order fulfillment: %w", err)
	}

	return &fulfillment, nil
}

// GetSellerOrders retrieves seller order summaries.
func (c *Client) GetSellerOrders(ctx context.Context, opts OrdersOptions) (*OrdersResponse, error) {
	params := buildOrdersParams(opts)
	resp, err := c.doRequest(ctx, "GET", "/seller/orders", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller orders: %w", err)
	}

	var orders OrdersResponse
	if err := c.decodeResponse(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode seller orders: %w", err)
	}

	return &orders, nil
}

// GetSellerOrder retrieves seller order details by ID.
func (c *Client) GetSellerOrder(ctx context.Context, id string) (*OrderDetailsResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/seller/orders/%s", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller order: %w", err)
	}

	var order OrderDetailsResponse
	if err := c.decodeResponse(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to decode seller order: %w", err)
	}

	return &order, nil
}

// UpdateSellerOrderFulfillment updates a seller order fulfillment.
func (c *Client) UpdateSellerOrderFulfillment(ctx context.Context, id string, req OrderFulfillmentRequest) (*OrderFulfillmentResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/seller/orders/%s/fulfillment", id)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update seller order fulfillment: %w", err)
	}

	var fulfillment OrderFulfillmentResponse
	if err := c.decodeResponse(resp, &fulfillment); err != nil {
		return nil, fmt.Errorf("failed to decode seller order fulfillment: %w", err)
	}

	return &fulfillment, nil
}

// GetSellerOrderReports retrieves order reports for a seller order.
func (c *Client) GetSellerOrderReports(ctx context.Context, id string) (*OrderReportsResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/seller/orders/%s/reports", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller order reports: %w", err)
	}

	var reports OrderReportsResponse
	if err := c.decodeResponse(resp, &reports); err != nil {
		return nil, fmt.Errorf("failed to decode seller order reports: %w", err)
	}

	return &reports, nil
}

func buildOrdersParams(opts OrdersOptions) url.Values {
	params := url.Values{}
	if opts.Since != nil {
		params.Add("since", opts.Since.Format(time.RFC3339Nano))
	}
	if opts.IsUnfulfilled != nil {
		params.Add("is_unfulfilled", strconv.FormatBool(*opts.IsUnfulfilled))
	}
	if opts.IsFulfilled != nil {
		params.Add("is_fulfilled", strconv.FormatBool(*opts.IsFulfilled))
	}
	if opts.HasFulfillments != nil {
		params.Add("has_fulfillments", strconv.FormatBool(*opts.HasFulfillments))
	}
	if opts.Label != "" {
		params.Add("label", opts.Label)
	}
	if opts.Limit > 0 {
		params.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Add("offset", strconv.Itoa(opts.Offset))
	}
	return params
}
