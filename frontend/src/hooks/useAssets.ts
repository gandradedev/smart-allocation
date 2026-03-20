import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import type { CreateAssetPayload, UpdateAssetPayload } from '../types/asset'

const QUERY_KEY = ['assets']

export function useAssets(totalToInvest?: number) {
  return useQuery({
    queryKey: [...QUERY_KEY, totalToInvest],
    queryFn: () => api.listAssets(totalToInvest),
  })
}

export function useCreateAsset() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateAssetPayload) => api.createAsset(payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })
}

export function useUpdateAsset() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ ticker, payload }: { ticker: string; payload: UpdateAssetPayload }) =>
      api.updateAsset(ticker, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })
}

export function useDeleteAsset() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ticker: string) => api.deleteAsset(ticker),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })
}
