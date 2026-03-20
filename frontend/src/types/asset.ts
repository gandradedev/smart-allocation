export interface Asset {
  ticker: string
  price: number
  quantity: number
  ceiling_price: number
  target_percent: number
  current_value: number
  current_percent: number
  target_value: number
  contribution_percent: number
  contribution_value: number
  shares_to_contribute: number
  ceiling_price_factor: number
  adjusted_contribution: number
}

export interface PortfolioResponse {
  assets: Asset[]
  total_assets: number
  total_value: number
}

export interface CreateAssetPayload {
  ticker: string
  price: number
  quantity: number
  ceiling_price: number
  target_percent: number
}

export type UpdateAssetPayload = Omit<CreateAssetPayload, 'ticker'>
