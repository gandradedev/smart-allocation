package asset

import (
	"context"

	"smart-allocation/internal/application/usecase/asset/dto"
	domainrepo "smart-allocation/internal/domain/repository"
)

//go:generate mockgen -destination=./mock/create_asset_use_case_mock.go -package=asset -source=create_asset_use_case.go CreateAssetUseCase
type CreateAssetUseCase interface {
	Execute(ctx context.Context, req *dto.CreateAssetRequestDTO) (*dto.CreateAssetResponseDTO, error)
}

type createAssetUseCase struct {
	repo domainrepo.AssetRepository
}

func NewCreateAssetUseCase(repo domainrepo.AssetRepository) CreateAssetUseCase {
	return &createAssetUseCase{repo: repo}
}

// Execute cria um novo ativo na carteira.
// Converte o DTO para entidade de domínio (aplicando validação) e persiste via repositório.
func (uc *createAssetUseCase) Execute(ctx context.Context, req *dto.CreateAssetRequestDTO) (*dto.CreateAssetResponseDTO, error) {
	asset, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return &dto.CreateAssetResponseDTO{Ticker: asset.Ticker}, nil
}
