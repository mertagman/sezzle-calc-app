export const MAX_DIGITS = 16
export const MAX_DECIMAL_PLACES = 12
export const SCIENTIFIC_UPPER_BOUND = 1e16
export const SCIENTIFIC_LOWER_BOUND = 1e-6

export function countDigits(str: string): number {
  return (str.match(/\d/g) || []).length
}

export function appendDigit(current: string, digit: string): string {
  if (countDigits(current) >= MAX_DIGITS) {
    return current
  }

  if (current === '0') {
    return digit
  }

  if (current === '-0') {
    return `-${digit}`
  }

  return current + digit
}

export function appendDecimal(current: string): string {
  if (current.includes('.')) {
    return current
  }

  if (countDigits(current) >= MAX_DIGITS) {
    return current
  }

  return `${current}.`
}

export function applyBackspace(current: string): string {
  if (
    current.length <= 1 ||
    current === '-0' ||
    (current.startsWith('-') && current.length === 2)
  ) {
    return '0'
  }

  const sliced = current.slice(0, -1)
  if (sliced === '' || sliced === '-') {
    return '0'
  }

  return sliced
}

export function toggleSign(current: string): string {
  if (current === '0') {
    return '0'
  }

  if (current.startsWith('-')) {
    return current.slice(1)
  }

  return `-${current}`
}

export function formatNumberResult(val: number): string {
  if (!Number.isFinite(val)) {
    return 'Error'
  }

  if (val === 0) {
    return '0'
  }

  const absVal = Math.abs(val)

  if (absVal >= SCIENTIFIC_UPPER_BOUND || absVal < SCIENTIFIC_LOWER_BOUND) {
    const expStr = Number(val.toPrecision(MAX_DECIMAL_PLACES)).toExponential()
    return expStr.replace(/e\+?/, 'e')
  }

  const precisionNum = Number(val.toPrecision(MAX_DECIMAL_PLACES))
  const precisionStr = precisionNum.toString()

  if (precisionStr.includes('.')) {
    const [, decimals] = precisionStr.split('.')
    if (decimals && decimals.length > MAX_DECIMAL_PLACES) {
      return Number(precisionNum.toFixed(MAX_DECIMAL_PLACES)).toString()
    }
  }

  return precisionStr
}
