package manapool

import (
	"context"
	"fmt"
)

// GetCardInfo retrieves card information for a list of card names.
func (c *Client) GetCardInfo(ctx context.Context, req CardInfoRequest) (*CardInfoResponse, error) {
	resp, err := c.doJSONRequest(ctx, "POST", "/card_info", nil, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get card info: %w", err)
	}

	var response CardInfoResponse
	if err := c.decodeResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to decode card info: %w", err)
	}

	return &response, nil
}
