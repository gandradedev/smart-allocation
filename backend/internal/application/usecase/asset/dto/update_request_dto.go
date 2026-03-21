package dto

import "smart-allocation/internal/domain/entity"

type UpdateAssetRequestDTO struct {
	Quantity      float64 `json:"quantity"`
	CeilingPrice  float64 `json:"ceiling_price"`
	TargetPercent float64 `json:"target_percent"`
}

func (req *UpdateAssetRequestDTO) ToEntity(ticker string) (*entity.Asset, error) {
	return entity.NewAsset(
		ticker,
		req.Quantity,
		req.CeilingPrice,
		req.TargetPercent,
	)
}
