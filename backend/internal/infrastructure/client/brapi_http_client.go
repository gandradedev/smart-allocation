package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	domainclient "smart-allocation/internal/domain/client"
)

const brapiBaseURL = "https://brapi.dev/api/quote"

type brapiResponse struct {
	Results []struct {
		Symbol             string  `json:"symbol"`
		RegularMarketPrice float64 `json:"regularMarketPrice"`
	} `json:"results"`
}

type brapiHTTPClient struct {
	token      string
	httpClient *http.Client
}

// NewBrapiHTTPClient creates a new HTTP client for brapi.dev.
func NewBrapiHTTPClient(token string) domainclient.BrapiClient {
	return &brapiHTTPClient{
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *brapiHTTPClient) GetQuote(ctx context.Context, ticker string) (*domainclient.AssetQuote, error) {
	url := fmt.Sprintf("%s/%s?token=%s", brapiBaseURL, ticker, c.token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("brapi: failed to build request for %s: %w", ticker, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("brapi: request failed for %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brapi: unexpected status %d for ticker %s", resp.StatusCode, ticker)
	}

	var body brapiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("brapi: failed to decode response for %s: %w", ticker, err)
	}

	if len(body.Results) == 0 {
		return nil, fmt.Errorf("brapi: no results returned for ticker %s", ticker)
	}

	r := body.Results[0]
	return &domainclient.AssetQuote{
		Symbol:             r.Symbol,
		RegularMarketPrice: r.RegularMarketPrice,
	}, nil
}
