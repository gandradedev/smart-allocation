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

export const fmt = {
  currency: (v: number) => currency.format(v),
  percent: (v: number) => `${percent.format(v)}%`,
  decimal: (v: number) => decimal.format(v),
}
