import { describe, expect, it } from 'vitest'
import {
  appendDecimal,
  appendDigit,
  applyBackspace,
  countDigits,
  formatNumberResult,
  toggleSign,
} from './format'

describe('countDigits', () => {
  it('counts only numeric digits', () => {
    expect(countDigits('12345')).toBe(5)
    expect(countDigits('-123.45')).toBe(5)
    expect(countDigits('0')).toBe(1)
    expect(countDigits('-0.00')).toBe(3)
  })
})

describe('appendDigit', () => {
  it('replaces single zero with new digit', () => {
    expect(appendDigit('0', '7')).toBe('7')
    expect(appendDigit('-0', '7')).toBe('-7')
  })

  it('appends digit when count is under 16', () => {
    expect(appendDigit('123', '4')).toBe('1234')
  })

  it('disallows entering more than 16 digits', () => {
    const sixteenDigits = '1234567890123456'
    expect(countDigits(sixteenDigits)).toBe(16)
    expect(appendDigit(sixteenDigits, '7')).toBe(sixteenDigits)

    const sixteenDigitsWithDec = '12345678.90123456'
    expect(countDigits(sixteenDigitsWithDec)).toBe(16)
    expect(appendDigit(sixteenDigitsWithDec, '9')).toBe(sixteenDigitsWithDec)
  })
})

describe('appendDecimal', () => {
  it('adds decimal point to whole number', () => {
    expect(appendDecimal('0')).toBe('0.')
    expect(appendDecimal('42')).toBe('42.')
  })

  it('prevents multiple decimal points', () => {
    expect(appendDecimal('42.')).toBe('42.')
    expect(appendDecimal('3.14')).toBe('3.14')
  })

  it('prevents adding decimal if already at 16 digits', () => {
    const sixteen = '1234567890123456'
    expect(appendDecimal(sixteen)).toBe(sixteen)
  })
})

describe('applyBackspace', () => {
  it('removes the last character', () => {
    expect(applyBackspace('123')).toBe('12')
    expect(applyBackspace('12.3')).toBe('12.')
    expect(applyBackspace('12.')).toBe('12')
  })

  it('resets to 0 when last remaining digit is deleted', () => {
    expect(applyBackspace('7')).toBe('0')
    expect(applyBackspace('-7')).toBe('0')
    expect(applyBackspace('0')).toBe('0')
  })
})

describe('toggleSign', () => {
  it('flips sign between positive and negative', () => {
    expect(toggleSign('42')).toBe('-42')
    expect(toggleSign('-42')).toBe('42')
    expect(toggleSign('3.14')).toBe('-3.14')
    expect(toggleSign('-3.14')).toBe('3.14')
  })

  it('does not negate zero', () => {
    expect(toggleSign('0')).toBe('0')
  })
})

describe('formatNumberResult', () => {
  it('formats zero and standard numbers cleanly', () => {
    expect(formatNumberResult(0)).toBe('0')
    expect(formatNumberResult(42)).toBe('42')
    expect(formatNumberResult(-42)).toBe('-42')
  })

  it('smoothes floating point precision inaccuracies', () => {
    expect(formatNumberResult(0.1 + 0.2)).toBe('0.3')
  })

  it('handles repeating fractions with max 12 precision', () => {
    const result = formatNumberResult(1 / 3)
    const decimals = result.split('.')[1] || ''
    expect(decimals.length).toBeLessThanOrEqual(12)
    expect(result.startsWith('0.33333333333')).toBe(true)
  })

  it('converts very large numbers to scientific notation', () => {
    const large = 1e16
    expect(formatNumberResult(large)).toContain('e')
  })

  it('converts very small numbers to scientific notation', () => {
    const small = 0.00000001
    expect(formatNumberResult(small)).toContain('e')
  })

  it('handles non-finite values safely', () => {
    expect(formatNumberResult(Infinity)).toBe('Error')
    expect(formatNumberResult(-Infinity)).toBe('Error')
    expect(formatNumberResult(NaN)).toBe('Error')
  })
})
