package asset

import (
	"context"
	"log"

	"smart-allocation/internal/application/usecase/asset/dto"
	domainclient "smart-allocation/internal/domain/client"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/create_asset_use_case_mock.go -package=asset -source=create_asset_use_case.go CreateAssetUseCase
type CreateAssetUseCase interface {
	Execute(ctx context.Context, req *dto.CreateAssetRequestDTO) (*dto.CreateAssetResponseDTO, error)
}

type createAssetUseCase struct {
	repo   domainrepo.AssetRepository
	client domainclient.BrapiClient
}

func NewCreateAssetUseCase(repo domainrepo.AssetRepository, client domainclient.BrapiClient) CreateAssetUseCase {
	return &createAssetUseCase{repo: repo, client: client}
}

// Execute cria um novo ativo na carteira com price=0 e dispara uma goroutine
// para buscar e persistir o preço, ícone e moeda via brapi.dev de forma assíncrona.
func (uc *createAssetUseCase) Execute(ctx context.Context, req *dto.CreateAssetRequestDTO) (*dto.CreateAssetResponseDTO, error) {
	asset, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	go func(ticker string) {
		quote, err := uc.client.GetQuote(context.Background(), ticker)
		if err != nil {
			log.Printf("async metadata fetch failed for %s: %v", ticker, err)
			return
		}
		if err := uc.repo.UpdateMetadata(context.Background(), ticker, quote.RegularMarketPrice, quote.LogoURL, quote.Currency); err != nil {
			log.Printf("async metadata update failed for %s: %v", ticker, err)
		}
	}(asset.Ticker)

	return &dto.CreateAssetResponseDTO{Ticker: asset.Ticker}, nil
}
