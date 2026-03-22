import { useState } from 'react'
import { fmt } from '../utils/format'
import type { Asset, AssetType } from '../types/asset'
import { ASSET_TYPES } from '../types/asset'

interface PortfolioSummaryProps {
  totalValue: number
  totalAssets: number
  assets: Asset[]
  totalToInvest: number | undefined
  onTotalToInvestChange: (value: number | undefined) => void
}

export function PortfolioSummary({
  totalValue,
  totalAssets,
  assets,
  totalToInvest,
  onTotalToInvestChange,
}: PortfolioSummaryProps) {
  const typeCounts = ASSET_TYPES.reduce<Partial<Record<AssetType, number>>>((acc, type) => {
    const count = assets.filter(a => a.asset_type === type).length
    if (count > 0) acc[type] = count
    return acc
  }, {})
  const [input, setInput] = useState(totalToInvest?.toString() ?? '')

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const raw = e.target.value
    setInput(raw)
    const parsed = parseFloat(raw.replace(',', '.'))
    onTotalToInvestChange(isNaN(parsed) || parsed <= 0 ? undefined : parsed)
  }

  return (
    <div className="mx-auto max-w-7xl px-6 py-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-medium uppercase tracking-wide text-slate-500">Portfolio Value</p>
          <p className="mt-1 text-2xl font-semibold text-slate-900">{fmt.currency(totalValue)}</p>
        </div>

        <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-medium uppercase tracking-wide text-slate-500">Assets</p>
          <div className="mt-1 flex items-center gap-3">
            <p className="text-2xl font-semibold text-slate-900">{totalAssets}</p>
            {totalAssets > 0 && (
              <div className="flex flex-wrap gap-1.5">
                {(Object.entries(typeCounts) as [AssetType, number][]).map(([type, count]) => (
                  <span key={type} className="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-600">
                    {type} <span className="font-semibold text-slate-800">{count}</span>
                  </span>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-medium uppercase tracking-wide text-slate-500">Amount to Invest</p>
          <div className="mt-1 flex items-center gap-2">
            <span className="text-sm text-slate-400">R$</span>
            <input
              type="number"
              min="0"
              step="0.01"
              placeholder="0,00"
              value={input}
              onChange={handleChange}
              className="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-1.5 text-lg font-semibold text-slate-900 placeholder-slate-300 focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            />
          </div>
          {totalToInvest && (
            <p className="mt-1 text-xs text-slate-400">Adjusted contributions calculated</p>
          )}
        </div>
      </div>
    </div>
  )
}
