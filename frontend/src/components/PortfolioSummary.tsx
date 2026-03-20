import { useState } from 'react'
import { fmt } from '../utils/format'

interface PortfolioSummaryProps {
  totalValue: number
  totalAssets: number
  totalToInvest: number | undefined
  onTotalToInvestChange: (value: number | undefined) => void
}

export function PortfolioSummary({
  totalValue,
  totalAssets,
  totalToInvest,
  onTotalToInvestChange,
}: PortfolioSummaryProps) {
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
          <p className="mt-1 text-2xl font-semibold text-slate-900">{totalAssets}</p>
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
