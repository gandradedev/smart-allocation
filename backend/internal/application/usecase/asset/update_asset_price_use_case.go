package asset

import (
	"context"

	domainclient "smart-allocation/internal/domain/client"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/update_asset_price_use_case_mock.go -package=asset -source=update_asset_price_use_case.go UpdateAssetPriceUseCase
type UpdateAssetPriceUseCase interface {
	Execute(ctx context.Context, ticker string) error
}

type updateAssetPriceUseCase struct {
	repo   domainrepo.AssetRepository
	client domainclient.BrapiClient
}

func NewUpdateAssetPriceUseCase(repo domainrepo.AssetRepository, client domainclient.BrapiClient) UpdateAssetPriceUseCase {
	return &updateAssetPriceUseCase{repo: repo, client: client}
}

// Execute fetches the current market price from brapi.dev and persists it.
func (uc *updateAssetPriceUseCase) Execute(ctx context.Context, ticker string) error {
	quote, err := uc.client.GetQuote(ctx, ticker)
	if err != nil {
		return err
	}

	return uc.repo.UpdatePrice(ctx, ticker, quote.RegularMarketPrice)
}
