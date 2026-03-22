import { useState } from 'react'
import type { Asset, AssetType } from '../types/asset'
import { ASSET_TYPES } from '../types/asset'
import { fmt } from '../utils/format'
import { usePollAssetPrice } from '../hooks/useAssets'

type SortKey = 'ticker' | 'current_value' | 'current_percent' | 'target_percent' | 'deviation'
type SortDir = 'asc' | 'desc'
type AllocationFilter = 'all' | 'under' | 'over' | 'on'

function getSortValue(asset: Asset, key: SortKey): number | string {
  if (key === 'ticker') return asset.ticker
  if (key === 'deviation') return asset.current_percent - asset.target_percent
  return asset[key]
}

function sortAssets(assets: Asset[], key: SortKey, dir: SortDir): Asset[] {
  return [...assets].sort((a, b) => {
    const av = getSortValue(a, key)
    const bv = getSortValue(b, key)
    const cmp = typeof av === 'string' ? av.localeCompare(bv as string) : (av as number) - (bv as number)
    return dir === 'asc' ? cmp : -cmp
  })
}

function filterByAllocation(assets: Asset[], filter: AllocationFilter): Asset[] {
  if (filter === 'all') return assets
  return assets.filter(a => {
    const diff = a.current_percent - a.target_percent
    if (filter === 'over') return diff > 0
    if (filter === 'under') return diff < 0
    return diff === 0
  })
}

function SortIcon({ active, dir }: { active: boolean; dir: SortDir }) {
  if (!active) return (
    <svg className="h-3 w-3 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4M17 8v12m0 0l4-4m-4 4l-4-4" />
    </svg>
  )
  return dir === 'asc'
    ? <svg className="h-3 w-3 text-slate-600" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" /></svg>
    : <svg className="h-3 w-3 text-slate-600" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
}

function PricePollCell({ ticker }: { ticker: string }) {
  usePollAssetPrice(ticker)
  return (
    <svg className="h-4 w-4 animate-spin text-slate-400" fill="none" viewBox="0 0 24 24">
      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
    </svg>
  )
}

interface AssetTableProps {
  assets: Asset[]
  totalToInvest: number | undefined
  onEdit: (asset: Asset) => void
  onDelete: (asset: Asset) => void
}

function AllocationBadge({ current, target }: { current: number; target: number }) {
  const diff = current - target
  if (diff === 0) {
    return <span className="inline-flex rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-600">On target</span>
  }
  if (diff > 0) {
    return <span className="inline-flex rounded-full bg-red-50 px-2 py-0.5 text-xs font-medium text-red-600">+{fmt.percent(diff)} over</span>
  }
  return <span className="inline-flex items-center gap-1 whitespace-nowrap rounded-full bg-green-50 px-2 py-0.5 text-xs font-medium text-green-600">{fmt.percent(Math.abs(diff))} under</span>
}

function SkeletonRow() {
  return (
    <tr>
      {Array.from({ length: 9 }).map((_, i) => (
        <td key={i} className="px-4 py-3">
          <div className="h-4 animate-pulse rounded bg-slate-200" />
        </td>
      ))}
    </tr>
  )
}

interface AssetTableWithLoadingProps extends AssetTableProps {
  isLoading: boolean
}

const ASSET_TYPE_LABELS: Record<AssetType, string> = {
  ACAO: 'Ações',
  FII: 'REITs (FII)',
  ETF: 'ETF',
  BDR: 'BDR',
  STOCK: 'Stocks',
}

export function AssetTable({ assets, totalToInvest, onEdit, onDelete, isLoading }: AssetTableWithLoadingProps) {
  const [sortKey, setSortKey] = useState<SortKey>('ticker')
  const [sortDir, setSortDir] = useState<SortDir>('asc')
  const [allocationFilter, setAllocationFilter] = useState<AllocationFilter>('all')

  function handleSort(key: SortKey) {
    if (key === sortKey) {
      setSortDir(d => d === 'asc' ? 'desc' : 'asc')
    } else {
      setSortKey(key)
      setSortDir('asc')
    }
  }

  function Th({ label, sortable, column }: { label: string; sortable?: SortKey; column?: string }) {
    if (!sortable) return <th className={`px-4 py-3 ${column ?? ''}`}>{label}</th>
    return (
      <th className={`px-4 py-3 ${column ?? ''}`}>
        <button
          onClick={() => handleSort(sortable)}
          className="inline-flex items-center gap-1 hover:text-slate-700"
        >
          {label}
          <SortIcon active={sortKey === sortable} dir={sortDir} />
        </button>
      </th>
    )
  }

  const FILTER_OPTIONS: { value: AllocationFilter; label: string }[] = [
    { value: 'all', label: 'All' },
    { value: 'under', label: 'Under' },
    { value: 'over', label: 'Over' },
    { value: 'on', label: 'On target' },
  ]

  if (!isLoading && assets.length === 0) {
    return (
      <div className="mx-auto max-w-7xl px-6 pb-12">
        <div className="rounded-xl border border-dashed border-slate-300 bg-white py-16 text-center">
          <svg className="mx-auto h-10 w-10 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          <p className="mt-3 text-sm font-medium text-slate-500">No assets yet</p>
          <p className="mt-1 text-xs text-slate-400">Add your first asset to start tracking your portfolio.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-7xl px-6 pb-12">
      <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        <div className="flex items-center gap-2 border-b border-slate-100 px-4 py-2">
          <span className="text-xs font-medium text-slate-400 uppercase tracking-wide">Filter</span>
          {FILTER_OPTIONS.map(opt => (
            <button
              key={opt.value}
              onClick={() => setAllocationFilter(opt.value)}
              className={`rounded-full px-3 py-1 text-xs font-medium transition-colors ${
                allocationFilter === opt.value
                  ? 'bg-slate-800 text-white'
                  : 'text-slate-500 hover:bg-slate-100'
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100 bg-slate-50 text-left text-xs font-medium tracking-wide text-slate-500">
                <Th label="Ticker" sortable="ticker" />
                <th className="px-4 py-3">Qty</th>
                <th className="px-4 py-3">Price</th>
                <th className="px-4 py-3">Ceiling Price</th>
                <Th label="Current Value" sortable="current_value" />
                <Th label="Current %" sortable="current_percent" />
                <Th label="Target %" sortable="target_percent" />
                <Th label="Allocation" sortable="deviation" />
                {totalToInvest && (
                  <>
                    <th className="px-4 py-3">Units to Buy</th>
                    <th className="px-4 py-3">Amount to Invest</th>
                  </>
                )}
                <th className="px-4 py-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {isLoading
                ? Array.from({ length: 4 }).map((_, i) => <SkeletonRow key={i} />)
                : ASSET_TYPES.flatMap(type => {
                    const rawGroup = assets.filter(a => a.asset_type === type)
                    const group = sortAssets(filterByAllocation(rawGroup, allocationFilter), sortKey, sortDir)
                    if (group.length === 0) return []

                    const groupValue = group.reduce((sum, a) => sum + a.current_value, 0)
                    const groupPercent = group.reduce((sum, a) => sum + a.current_percent, 0)
                    // 9 base columns (Ticker, Qty, Price, Ceiling, Value, Current%, Target%, Allocation, Actions)
                    // + 2 optional columns when totalToInvest is set (Units to Buy, Amount to Invest)
                    const colSpan = 9 + (totalToInvest ? 2 : 0)

                    return [
                      <tr key={`group-${type}`} className="bg-slate-100">
                        <td colSpan={colSpan} className="px-4 py-2">
                          <div className="flex items-center gap-3">
                            <span className="text-xs font-semibold uppercase tracking-wide text-slate-600">
                              {ASSET_TYPE_LABELS[type]}
                            </span>
                            <span className="text-xs text-slate-400">{group.length} {group.length === 1 ? 'asset' : 'assets'}</span>
                            <span className="ml-auto text-xs font-medium text-slate-600">{fmt.currency(groupValue)}</span>
                            <span className="text-xs font-medium text-slate-600">{fmt.percent(groupPercent)}</span>
                          </div>
                        </td>
                      </tr>,
                      ...group.map(asset => (
                        <tr key={asset.ticker} className="hover:bg-slate-50">
                          <td className="px-4 py-3">
                            <div className="flex items-center gap-2">
                              {asset.icon && asset.asset_type === 'ACAO' ? (
                                <img
                                  src={asset.icon}
                                  alt={asset.ticker}
                                  className="h-6 w-6 rounded-full object-contain"
                                />
                              ) : (
                                <div className="flex h-6 w-6 items-center justify-center rounded-full bg-slate-200 text-[10px] font-bold text-slate-500">
                                  {asset.ticker.slice(0, 2)}
                                </div>
                              )}
                              <span className="font-semibold text-slate-900">{asset.ticker}</span>
                            </div>
                          </td>
                          <td className="px-4 py-3 text-slate-700">{fmt.decimal(asset.quantity)}</td>
                          <td className="px-4 py-3 text-slate-700">
                            {asset.price === 0
                              ? <PricePollCell ticker={asset.ticker} />
                              : fmt.currencyByCode(asset.price, asset.currency)
                            }
                          </td>
                          <td className="px-4 py-3 text-slate-700">{fmt.currencyByCode(asset.ceiling_price, asset.currency)}</td>
                          <td className="px-4 py-3 text-slate-700">{fmt.currencyByCode(asset.current_value, asset.currency)}</td>
                          <td className="px-4 py-3 text-slate-700">{fmt.percent(asset.current_percent)}</td>
                          <td className="px-4 py-3 text-slate-700">{fmt.percent(asset.target_percent)}</td>
                          <td className="px-4 py-3">
                            <AllocationBadge current={asset.current_percent} target={asset.target_percent} />
                          </td>
                          {totalToInvest && (
                            <>
                              <td className={`px-4 py-3 font-medium ${asset.shares_to_contribute > 0 ? 'text-green-700' : 'text-slate-700'}`}>
                                {asset.shares_to_contribute > 0 ? fmt.decimal(asset.shares_to_contribute) : '—'}
                              </td>
                              <td className={`px-4 py-3 font-medium ${asset.adjusted_contribution > 0 ? 'text-green-700' : 'text-slate-700'}`}>
                                {asset.adjusted_contribution > 0 ? fmt.currencyByCode(asset.adjusted_contribution, asset.currency) : '—'}
                              </td>
                            </>
                          )}
                          <td className="px-4 py-3 text-right">
                            <div className="flex justify-end gap-1">
                              <button
                                onClick={() => onEdit(asset)}
                                className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-100 hover:text-slate-700"
                                title="Edit"
                              >
                                <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                </svg>
                              </button>
                              <button
                                onClick={() => onDelete(asset)}
                                className="rounded-lg p-1.5 text-slate-400 hover:bg-red-50 hover:text-red-600"
                                title="Delete"
                              >
                                <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                </svg>
                              </button>
                            </div>
                          </td>
                        </tr>
                      )),
                    ]
                  })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
