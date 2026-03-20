package asset

import (
	"context"

	"smart-allocation/internal/application/usecase/asset/dto"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/update_asset_use_case_mock.go -package=asset -source=update_asset_use_case.go UpdateAssetUseCase
type UpdateAssetUseCase interface {
	Execute(ctx context.Context, ticker string, req *dto.UpdateAssetRequestDTO) error
}

type updateAssetUseCase struct {
	repo domainrepo.AssetRepository
}

func NewUpdateAssetUseCase(repo domainrepo.AssetRepository) UpdateAssetUseCase {
	return &updateAssetUseCase{repo: repo}
}

// Execute atualiza os dados de um ativo existente.
func (uc *updateAssetUseCase) Execute(ctx context.Context, ticker string, req *dto.UpdateAssetRequestDTO) error {
	asset, err := req.ToEntity(ticker)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, ticker, asset)
}
