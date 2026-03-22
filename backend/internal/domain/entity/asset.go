package entity

import (
	"slices"
	"strings"

	domainerrors "smart-allocation/internal/domain/errors"
)

type AssetType string

const (
	AssetTypeAcao  AssetType = "ACAO"
	AssetTypeFII   AssetType = "FII"
	AssetTypeETF   AssetType = "ETF"
	AssetTypeBDR   AssetType = "BDR"
	AssetTypeStock AssetType = "STOCK"
)

var validAssetTypes = []AssetType{AssetTypeAcao, AssetTypeFII, AssetTypeETF, AssetTypeBDR, AssetTypeStock}

func (t AssetType) isValid() bool {
	return slices.Contains(validAssetTypes, t)
}

// Asset is the core domain entity. It has no JSON tags as it is not exposed directly via HTTP.
type Asset struct {
	Ticker        string
	AssetType     AssetType
	Quantity      float64
	Price         float64
	CeilingPrice  float64
	TargetPercent float64
	Icon          string
	Currency      string
}

// NewAsset creates and validates a new asset.
func NewAsset(ticker string, assetType AssetType, quantity, ceilingPrice, targetPercent float64) (*Asset, error) {
	a := &Asset{
		Ticker:        strings.ToUpper(strings.TrimSpace(ticker)),
		AssetType:     assetType,
		Quantity:      quantity,
		CeilingPrice:  ceilingPrice,
		TargetPercent: targetPercent,
	}

	if err := a.Validate(); err != nil {
		return nil, err
	}

	return a, nil
}

// Validate checks whether the asset has all required fields.
func (a *Asset) Validate() error {
	if a.Ticker == "" {
		return domainerrors.NewValidationError("ticker_required", "Ticker is required", nil)
	}
	if !a.AssetType.isValid() {
		return domainerrors.NewValidationError("invalid_asset_type", "Asset type must be one of: ACAO, FII, ETF, BDR, STOCK", nil)
	}
	return nil
}

// CeilingPriceFactor calculates the attractiveness factor of the asset for rebalancing.
// Returns 0 if the current price has reached or exceeded the ceiling price,
// or if the current percent has reached or exceeded the target.
func (a *Asset) CeilingPriceFactor(currentPercent float64) float64 {
	if a.Price >= a.CeilingPrice || currentPercent >= a.TargetPercent {
		return 0
	}
	return ((a.CeilingPrice - a.Price) / a.CeilingPrice) * (a.TargetPercent - currentPercent)
}
