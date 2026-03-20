package asset

import (
	"context"
	"log"

	"smart-allocation/internal/application/usecase/asset/dto"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/create_asset_use_case_mock.go -package=asset -source=create_asset_use_case.go CreateAssetUseCase
type CreateAssetUseCase interface {
	Execute(ctx context.Context, req *dto.CreateAssetRequestDTO) (*dto.CreateAssetResponseDTO, error)
}

type createAssetUseCase struct {
	repo         domainrepo.AssetRepository
	priceUpdater UpdateAssetPriceUseCase
}

func NewCreateAssetUseCase(repo domainrepo.AssetRepository, priceUpdater UpdateAssetPriceUseCase) CreateAssetUseCase {
	return &createAssetUseCase{repo: repo, priceUpdater: priceUpdater}
}

// Execute cria um novo ativo na carteira com price=0 e dispara uma goroutine
// para buscar e persistir o preço atual via brapi.dev de forma assíncrona.
func (uc *createAssetUseCase) Execute(ctx context.Context, req *dto.CreateAssetRequestDTO) (*dto.CreateAssetResponseDTO, error) {
	asset, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	go func(ticker string) {
		if err := uc.priceUpdater.Execute(context.Background(), ticker); err != nil {
			log.Printf("async price update failed for %s: %v", ticker, err)
		}
	}(asset.Ticker)

	return &dto.CreateAssetResponseDTO{Ticker: asset.Ticker}, nil
}
