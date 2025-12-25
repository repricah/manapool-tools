package manapool

import (
	"context"
	"fmt"
	"net/url"
)

// GetWebhooks retrieves registered webhooks.
func (c *Client) GetWebhooks(ctx context.Context, topic string) (*WebhooksResponse, error) {
	params := url.Values{}
	if topic != "" {
		params.Add("topic", topic)
	}

	resp, err := c.doRequest(ctx, "GET", "/webhooks", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks: %w", err)
	}

	var webhooks WebhooksResponse
	if err := c.decodeResponse(resp, &webhooks); err != nil {
		return nil, fmt.Errorf("failed to decode webhooks: %w", err)
	}

	return &webhooks, nil
}

// GetWebhook retrieves a webhook by ID.
func (c *Client) GetWebhook(ctx context.Context, id string) (*Webhook, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/webhooks/%s", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	var webhook Webhook
	if err := c.decodeResponse(resp, &webhook); err != nil {
		return nil, fmt.Errorf("failed to decode webhook: %w", err)
	}

	return &webhook, nil
}

// RegisterWebhook registers a webhook.
func (c *Client) RegisterWebhook(ctx context.Context, req WebhookRegisterRequest) (*Webhook, error) {
	resp, err := c.doJSONRequest(ctx, "PUT", "/webhooks/register", nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register webhook: %w", err)
	}

	var webhook Webhook
	if err := c.decodeResponse(resp, &webhook); err != nil {
		return nil, fmt.Errorf("failed to decode webhook registration: %w", err)
	}

	return &webhook, nil
}

// DeleteWebhook deletes a webhook by ID.
func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	if id == "" {
		return NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/webhooks/%s", id)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return c.decodeResponse(resp, nil)
}
