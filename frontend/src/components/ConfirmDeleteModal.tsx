import { Modal } from './ui/Modal'

interface ConfirmDeleteModalProps {
  ticker: string
  isPending: boolean
  onConfirm: () => void
  onClose: () => void
}

export function ConfirmDeleteModal({ ticker, isPending, onConfirm, onClose }: ConfirmDeleteModalProps) {
  return (
    <Modal title="Remove Asset" onClose={onClose}>
      <div className="space-y-5">
        <p className="text-sm text-slate-600">
          Are you sure you want to remove <span className="font-semibold text-slate-900">{ticker}</span> from your portfolio? This action cannot be undone.
        </p>
        <div className="flex justify-end gap-3">
          <button
            onClick={onClose}
            className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          >
            Cancel
          </button>
          <button
            onClick={onConfirm}
            disabled={isPending}
            className="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-60"
          >
            {isPending ? 'Removing…' : 'Remove'}
          </button>
        </div>
      </div>
    </Modal>
  )
}
