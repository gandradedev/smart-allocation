const currency = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' })
const percent = new Intl.NumberFormat('pt-BR', {
  style: 'decimal',
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})
const decimal = new Intl.NumberFormat('pt-BR', {
  minimumFractionDigits: 0,
  maximumFractionDigits: 2,
})

function currencyByCode(v: number, code: string): string {
  if (!code || code === 'BRL') return currency.format(v)
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: code }).format(v)
}

export const fmt = {
  currency: (v: number) => currency.format(v),
  currencyByCode,
  percent: (v: number) => `${percent.format(v)}%`,
  decimal: (v: number) => decimal.format(v),
}
