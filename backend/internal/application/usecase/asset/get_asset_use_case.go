package asset

import (
	"context"

	"smart-allocation/internal/application/usecase/asset/dto"
	domainerrors "smart-allocation/internal/domain/errors"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/get_asset_use_case_mock.go -package=asset -source=get_asset_use_case.go GetAssetUseCase
type GetAssetUseCase interface {
	Execute(ctx context.Context, ticker string, totalToInvest float64) (*dto.GetAssetResponseDTO, error)
}

type getAssetUseCase struct {
	repo domainrepo.AssetRepository
}

func NewGetAssetUseCase(repo domainrepo.AssetRepository) GetAssetUseCase {
	return &getAssetUseCase{repo: repo}
}

// Execute fetches an asset by ticker and returns its data with calculated rebalancing fields.
// Uses FindAll to compute totalValue and sumFactors in a single query,
// avoiding extra DB calls to get full portfolio context.
func (uc *getAssetUseCase) Execute(ctx context.Context, ticker string, totalToInvest float64) (*dto.GetAssetResponseDTO, error) {
	assets, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var totalValue float64
	for _, a := range assets {
		totalValue += a.Quantity * a.Price
	}

	var sumFactors float64
	for _, a := range assets {
		var currentPercent float64
		if totalValue > 0 {
			currentPercent = (a.Quantity * a.Price / totalValue) * 100
		}
		sumFactors += a.CeilingPriceFactor(currentPercent)
	}

	for _, a := range assets {
		if a.Ticker == ticker {
			return (&dto.GetAssetResponseDTO{}).FromEntity(a, totalValue, sumFactors, totalToInvest), nil
		}
	}

	return nil, domainerrors.NewNotFoundError("Asset not found")
}
