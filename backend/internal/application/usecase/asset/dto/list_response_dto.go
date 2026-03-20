package dto

// ListAssetsResponseDTO is the output DTO for the full portfolio listing.
type ListAssetsResponseDTO struct {
	TotalValue  float64               `json:"total_value"`
	TotalAssets int                   `json:"total_assets"`
	Assets      []GetAssetResponseDTO `json:"assets"`
}
