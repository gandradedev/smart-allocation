package repository

import (
	"context"

	"smart-allocation/internal/domain/entity"
)

// AssetRepository defines the persistence contract for assets.
// The interface lives in the domain to ensure infrastructure depends on the domain, not the other way around.
type AssetRepository interface {
	Create(ctx context.Context, asset *entity.Asset) error
	FindAll(ctx context.Context) ([]*entity.Asset, error)
	FindByTicker(ctx context.Context, ticker string) (*entity.Asset, error)
	TotalValue(ctx context.Context) (float64, error)
	Update(ctx context.Context, ticker string, asset *entity.Asset) error
	UpdatePrice(ctx context.Context, ticker string, price float64) error
	UpdateMetadata(ctx context.Context, ticker string, price float64, icon, currency string) error
	Delete(ctx context.Context, ticker string) error
}
