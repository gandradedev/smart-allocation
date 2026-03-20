package asset

import (
	"context"

	"smart-allocation/internal/application/usecase/asset/dto"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/list_assets_use_case_mock.go -package=asset -source=list_assets_use_case.go ListAssetsUseCase
type ListAssetsUseCase interface {
	Execute(ctx context.Context, totalToInvest float64) (*dto.ListAssetsResponseDTO, error)
}

type listAssetsUseCase struct {
	repo domainrepo.AssetRepository
}

func NewListAssetsUseCase(repo domainrepo.AssetRepository) ListAssetsUseCase {
	return &listAssetsUseCase{repo: repo}
}

// Execute returns all assets in the portfolio with the rebalancing summary.
// totalToInvest is the total amount available for contribution in this cycle (spreadsheet param $B$19).
func (uc *listAssetsUseCase) Execute(ctx context.Context, totalToInvest float64) (*dto.ListAssetsResponseDTO, error) {
	assets, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var totalValue float64
	for _, a := range assets {
		totalValue += a.Quantity * a.Price
	}

	// Computes the sum of all factors (needed to distribute the contribution proportionally).
	// Corresponds to $E$18 in the spreadsheet.
	var sumFactors float64
	for _, a := range assets {
		var currentPercent float64
		if totalValue > 0 {
			currentPercent = (a.Quantity * a.Price / totalValue) * 100
		}
		sumFactors += a.CeilingPriceFactor(currentPercent)
	}

	responses := make([]dto.GetAssetResponseDTO, len(assets))
	for i, a := range assets {
		responses[i] = *(&dto.GetAssetResponseDTO{}).FromEntity(a, totalValue, sumFactors, totalToInvest)
	}

	return &dto.ListAssetsResponseDTO{
		TotalValue:  totalValue,
		TotalAssets: len(assets),
		Assets:      responses,
	}, nil
}
