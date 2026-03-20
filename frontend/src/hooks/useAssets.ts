import { useEffect } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import type { CreateAssetPayload, UpdateAssetPayload } from '../types/asset'

export const ASSETS_QUERY_KEY = ['assets']

export function useAssets(totalToInvest?: number) {
  return useQuery({
    queryKey: [...ASSETS_QUERY_KEY, totalToInvest],
    queryFn: () => api.listAssets(totalToInvest),
  })
}

export function usePollAssetPrice(ticker: string) {
  const queryClient = useQueryClient()

  const query = useQuery({
    queryKey: ['asset-price-poll', ticker],
    queryFn: () => api.getAsset(ticker),
    refetchInterval: query => (query.state.data?.price ?? 0) === 0 ? 2000 : false,
  })

  useEffect(() => {
    if (query.data && query.data.price > 0) {
      queryClient.invalidateQueries({ queryKey: ASSETS_QUERY_KEY })
    }
  }, [query.data?.price, queryClient])
}

export function useCreateAsset() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateAssetPayload) => api.createAsset(payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ASSETS_QUERY_KEY }),
  })
}

export function useUpdateAsset() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ ticker, payload }: { ticker: string; payload: UpdateAssetPayload }) =>
      api.updateAsset(ticker, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ASSETS_QUERY_KEY }),
  })
}

export function useDeleteAsset() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ticker: string) => api.deleteAsset(ticker),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ASSETS_QUERY_KEY }),
  })
}
