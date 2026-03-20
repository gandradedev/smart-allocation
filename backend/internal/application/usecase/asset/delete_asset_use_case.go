package asset

import (
	"context"

	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/delete_asset_use_case_mock.go -package=asset -source=delete_asset_use_case.go DeleteAssetUseCase
type DeleteAssetUseCase interface {
	Execute(ctx context.Context, ticker string) error
}

type deleteAssetUseCase struct {
	repo domainrepo.AssetRepository
}

func NewDeleteAssetUseCase(repo domainrepo.AssetRepository) DeleteAssetUseCase {
	return &deleteAssetUseCase{repo: repo}
}

// Execute remove um ativo da carteira pelo ticker.
func (uc *deleteAssetUseCase) Execute(ctx context.Context, ticker string) error {
	return uc.repo.Delete(ctx, ticker)
}
