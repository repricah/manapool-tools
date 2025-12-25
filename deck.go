package manapool

import (
	"context"
	"fmt"
)

// CreateDeck validates a deck and returns details.
func (c *Client) CreateDeck(ctx context.Context, req DeckCreateRequest) (*DeckCreateResponse, error) {
	resp, err := c.doJSONRequest(ctx, "POST", "/deck", nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	var response DeckCreateResponse
	if err := c.decodeResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to decode deck response: %w", err)
	}

	return &response, nil
}
