package dto

import "smart-allocation/internal/domain/entity"

type CreateAssetRequestDTO struct {
	Ticker        string  `json:"ticker"`
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	CeilingPrice  float64 `json:"ceiling_price"`
	TargetPercent float64 `json:"target_percent"`
}

func (req *CreateAssetRequestDTO) ToEntity() (*entity.Asset, error) {
	return entity.NewAsset(
		req.Ticker,
		req.Quantity,
		req.Price,
		req.CeilingPrice,
		req.TargetPercent,
	)
}
