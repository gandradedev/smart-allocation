import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Modal } from './ui/Modal'
import type { Asset, CreateAssetPayload, UpdateAssetPayload } from '../types/asset'

const schema = z.object({
  ticker: z.string().min(1, 'Required').max(10).transform(v => v.toUpperCase()),
  price: z.coerce.number().positive('Must be greater than 0'),
  quantity: z.coerce.number().min(0, 'Cannot be negative'),
  ceiling_price: z.coerce.number().positive('Must be greater than 0'),
  target_percent: z.coerce.number().positive('Must be greater than 0').max(100, 'Max 100%'),
})

type FormValues = z.infer<typeof schema>

interface AssetFormModalProps {
  asset: Asset | null
  isPending: boolean
  onSubmit: (data: CreateAssetPayload | UpdateAssetPayload, ticker?: string) => void
  onClose: () => void
}

interface FieldProps {
  label: string
  error?: string
  children: React.ReactNode
}

function Field({ label, error, children }: FieldProps) {
  return (
    <div>
      <label className="mb-1 block text-sm font-medium text-slate-700">{label}</label>
      {children}
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  )
}

const inputClass =
  'w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-900 placeholder-slate-400 focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-500/20 disabled:opacity-50'

export function AssetFormModal({ asset, isPending, onSubmit, onClose }: AssetFormModalProps) {
  const isEditing = asset !== null

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: asset
      ? {
          ticker: asset.ticker,
          price: asset.price,
          quantity: asset.quantity,
          ceiling_price: asset.ceiling_price,
          target_percent: asset.target_percent,
        }
      : undefined,
  })

  useEffect(() => {
    if (asset) {
      reset({
        ticker: asset.ticker,
        price: asset.price,
        quantity: asset.quantity,
        ceiling_price: asset.ceiling_price,
        target_percent: asset.target_percent,
      })
    } else {
      reset({})
    }
  }, [asset, reset])

  function onValid(data: FormValues) {
    if (isEditing) {
      const { ticker: _ticker, ...payload } = data
      onSubmit(payload as UpdateAssetPayload, asset.ticker)
    } else {
      onSubmit(data as CreateAssetPayload)
    }
  }

  return (
    <Modal title={isEditing ? `Edit ${asset.ticker}` : 'New Asset'} onClose={onClose}>
      <form onSubmit={handleSubmit(onValid)} className="space-y-4">
        <Field label="Ticker" error={errors.ticker?.message}>
          <input
            {...register('ticker')}
            placeholder="BBAS3"
            disabled={isEditing}
            className={inputClass}
          />
        </Field>

        <div className="grid grid-cols-2 gap-4">
          <Field label="Price (R$)" error={errors.price?.message}>
            <input {...register('price')} type="number" step="0.01" placeholder="0.00" className={inputClass} />
          </Field>

          <Field label="Quantity" error={errors.quantity?.message}>
            <input {...register('quantity')} type="number" step="0.01" placeholder="0" className={inputClass} />
          </Field>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Field label="Ceiling Price (R$)" error={errors.ceiling_price?.message}>
            <input {...register('ceiling_price')} type="number" step="0.01" placeholder="0.00" className={inputClass} />
          </Field>

          <Field label="Target %" error={errors.target_percent?.message}>
            <input {...register('target_percent')} type="number" step="0.01" placeholder="0.00" className={inputClass} />
          </Field>
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isPending}
            className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-60"
          >
            {isPending ? 'Saving…' : isEditing ? 'Save Changes' : 'Add Asset'}
          </button>
        </div>
      </form>
    </Modal>
  )
}
