package client

import "context"

// AssetQuote holds the price data returned by the brapi.dev API.
type AssetQuote struct {
	Symbol             string
	RegularMarketPrice float64
	LogoURL            string
	Currency           string
}

// BrapiClient is the domain contract for fetching asset quotes.
// The interface lives in the domain so the application layer depends on the domain, not on infrastructure.
type BrapiClient interface {
	GetQuote(ctx context.Context, ticker string) (*AssetQuote, error)
}
