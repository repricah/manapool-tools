package manapool

import (
	"context"
	"fmt"
)

// GetSinglesPrices retrieves prices for all in-stock singles.
func (c *Client) GetSinglesPrices(ctx context.Context) (*SinglesPricesList, error) {
	resp, err := c.doRequest(ctx, "GET", "/prices/singles", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get singles prices: %w", err)
	}

	var prices SinglesPricesList
	if err := c.decodeResponse(resp, &prices); err != nil {
		return nil, fmt.Errorf("failed to decode singles prices: %w", err)
	}

	return &prices, nil
}

// GetVariantPrices retrieves prices for all in-stock variants.
func (c *Client) GetVariantPrices(ctx context.Context) (*VariantPricesList, error) {
	resp, err := c.doRequest(ctx, "GET", "/prices/variants", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get variant prices: %w", err)
	}

	var prices VariantPricesList
	if err := c.decodeResponse(resp, &prices); err != nil {
		return nil, fmt.Errorf("failed to decode variant prices: %w", err)
	}

	return &prices, nil
}

// GetSealedPrices retrieves prices for all in-stock sealed products.
func (c *Client) GetSealedPrices(ctx context.Context) (*SealedPricesList, error) {
	resp, err := c.doRequest(ctx, "GET", "/prices/sealed", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get sealed prices: %w", err)
	}

	var prices SealedPricesList
	if err := c.decodeResponse(resp, &prices); err != nil {
		return nil, fmt.Errorf("failed to decode sealed prices: %w", err)
	}

	return &prices, nil
}
