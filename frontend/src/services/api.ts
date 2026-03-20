import type { CreateAssetPayload, PortfolioResponse, UpdateAssetPayload } from '../types/asset'

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: 'Unexpected error' }))
    throw new Error(error.message ?? `HTTP ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export const api = {
  listAssets(totalToInvest?: number): Promise<PortfolioResponse> {
    const params = totalToInvest ? `?total_to_invest=${totalToInvest}` : ''
    return fetch(`/assets${params}`).then(handleResponse<PortfolioResponse>)
  },

  createAsset(payload: CreateAssetPayload): Promise<{ ticker: string }> {
    return fetch('/assets', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }).then(handleResponse<{ ticker: string }>)
  },

  updateAsset(ticker: string, payload: UpdateAssetPayload): Promise<void> {
    return fetch(`/assets/${ticker}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }).then(handleResponse<void>)
  },

  deleteAsset(ticker: string): Promise<void> {
    return fetch(`/assets/${ticker}`, { method: 'DELETE' }).then(handleResponse<void>)
  },
}
