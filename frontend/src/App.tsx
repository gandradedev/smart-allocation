import { useState } from 'react'
import { Header } from './components/Header'
import { PortfolioSummary } from './components/PortfolioSummary'
import { AssetTable } from './components/AssetTable'
import { AssetFormModal } from './components/AssetFormModal'
import { ConfirmDeleteModal } from './components/ConfirmDeleteModal'
import { useAssets, useCreateAsset, useUpdateAsset, useDeleteAsset } from './hooks/useAssets'
import type { Asset, CreateAssetPayload, UpdateAssetPayload } from './types/asset'

export default function App() {
  const [totalToInvest, setTotalToInvest] = useState<number | undefined>(() => {
    const stored = localStorage.getItem('totalToInvest')
    return stored ? Number(stored) : undefined
  })

  function handleTotalToInvestChange(value: number | undefined) {
    if (value === undefined) {
      localStorage.removeItem('totalToInvest')
    } else {
      localStorage.setItem('totalToInvest', String(value))
    }
    setTotalToInvest(value)
  }
  const [assetToEdit, setAssetToEdit] = useState<Asset | null>(null)
  const [assetToDelete, setAssetToDelete] = useState<Asset | null>(null)
  const [isFormOpen, setIsFormOpen] = useState(false)

  const { data, isLoading, isError } = useAssets(totalToInvest)
  const createAsset = useCreateAsset()
  const updateAsset = useUpdateAsset()
  const deleteAsset = useDeleteAsset()

  function openCreate() {
    setAssetToEdit(null)
    setIsFormOpen(true)
  }

  function openEdit(asset: Asset) {
    setAssetToEdit(asset)
    setIsFormOpen(true)
  }

  function closeForm() {
    setIsFormOpen(false)
    setAssetToEdit(null)
  }

  function handleFormSubmit(payload: CreateAssetPayload | UpdateAssetPayload, ticker?: string) {
    if (ticker) {
      updateAsset.mutate(
        { ticker, payload: payload as UpdateAssetPayload },
        { onSuccess: closeForm },
      )
    } else {
      createAsset.mutate(payload as CreateAssetPayload, { onSuccess: closeForm })
    }
  }

  function handleDelete() {
    if (!assetToDelete) return
    deleteAsset.mutate(assetToDelete.ticker, {
      onSuccess: () => setAssetToDelete(null),
    })
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <Header onNewAsset={openCreate} />

      {isError ? (
        <div className="mx-auto max-w-7xl px-6 py-12 text-center">
          <p className="text-sm text-red-500">Failed to connect to the backend. Make sure the server is running at <code className="rounded bg-slate-100 px-1">http://localhost:8080</code>.</p>
        </div>
      ) : (
        <>
          <PortfolioSummary
            totalValue={data?.total_value ?? 0}
            totalAssets={data?.total_assets ?? 0}
            totalToInvest={totalToInvest}
            onTotalToInvestChange={handleTotalToInvestChange}
          />

          <AssetTable
            assets={data?.assets ?? []}
            totalToInvest={totalToInvest}
            isLoading={isLoading}
            onEdit={openEdit}
            onDelete={setAssetToDelete}
          />
        </>
      )}

      {isFormOpen && (
        <AssetFormModal
          asset={assetToEdit}
          isPending={createAsset.isPending || updateAsset.isPending}
          onSubmit={handleFormSubmit}
          onClose={closeForm}
        />
      )}

      {assetToDelete && (
        <ConfirmDeleteModal
          ticker={assetToDelete.ticker}
          isPending={deleteAsset.isPending}
          onConfirm={handleDelete}
          onClose={() => setAssetToDelete(null)}
        />
      )}
    </div>
  )
}
