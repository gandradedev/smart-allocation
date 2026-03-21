package dto

import (
	"math"

	"smart-allocation/internal/domain/entity"
)

// GetAssetResponseDTO is the output DTO for a single asset, including calculated rebalancing fields.
type GetAssetResponseDTO struct {
	Ticker        string           `json:"ticker"`
	AssetType     entity.AssetType `json:"asset_type"`
	Quantity      float64          `json:"quantity"`
	Price         float64          `json:"price"`
	CeilingPrice  float64          `json:"ceiling_price"`
	TargetPercent float64          `json:"target_percent"`
	// Calculated fields
	TargetValue          float64 `json:"target_value"`
	CurrentPercent       float64 `json:"current_percent"`
	CurrentValue         float64 `json:"current_value"`
	ContributionPercent  float64 `json:"contribution_percent"`
	ContributionValue    float64 `json:"contribution_value"`
	SharesToContribute   float64 `json:"shares_to_contribute"`
	CeilingPriceFactor   float64 `json:"ceiling_price_factor"`
	AdjustedContribution float64 `json:"adjusted_contribution"`
}

// FromEntity converts the domain entity to the response DTO, applying rebalancing formulas.
//
// Parameters:
//   - totalValue: total portfolio value (sum of quantity * price for all assets)
//   - sumFactors: sum of ceiling_price_factor for all assets in the portfolio
//   - totalToInvest: total amount available to invest in this cycle
func (d *GetAssetResponseDTO) FromEntity(a *entity.Asset, totalValue, sumFactors, totalToInvest float64) *GetAssetResponseDTO {
	currentValue := a.Quantity * a.Price

	var currentPercent float64
	if totalValue > 0 {
		currentPercent = (currentValue / totalValue) * 100
	}

	targetValue := totalValue * (a.TargetPercent / 100)
	contributionPercent := a.TargetPercent - currentPercent
	contributionValue := targetValue - currentValue

	var sharesToContribute float64
	if a.Price > 0 && contributionPercent > 0 {
		sharesToContribute = math.Ceil(contributionValue / a.Price)
	}

	ceilingPriceFactor := a.CeilingPriceFactor(currentPercent)

	var adjustedContribution float64
	if sumFactors > 0 {
		adjustedContribution = totalToInvest * (ceilingPriceFactor / sumFactors)
	}

	return &GetAssetResponseDTO{
		Ticker:               a.Ticker,
		AssetType:            a.AssetType,
		Quantity:             a.Quantity,
		Price:                a.Price,
		CeilingPrice:         a.CeilingPrice,
		TargetPercent:        a.TargetPercent,
		TargetValue:          targetValue,
		CurrentPercent:       currentPercent,
		CurrentValue:         currentValue,
		ContributionPercent:  contributionPercent,
		ContributionValue:    contributionValue,
		SharesToContribute:   sharesToContribute,
		CeilingPriceFactor:   ceilingPriceFactor,
		AdjustedContribution: adjustedContribution,
	}
}
