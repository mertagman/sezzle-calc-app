import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import * as client from '../api/client'
import { useCalculator } from './useCalculator'

vi.mock('../api/client', () => ({
  calculateBinary: vi.fn(),
  calculateUnary: vi.fn(),
  checkHealth: vi.fn(),
}))

describe('useCalculator', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    localStorage.clear()
  })

  it('initializes with default values', () => {
    const { result } = renderHook(() => useCalculator())
    expect(result.current.displayValue).toBe('0')
    expect(result.current.subDisplayValue).toBe('')
    expect(result.current.error).toBeNull()
    expect(result.current.history).toEqual([])
  })

  it('inputs digits and prevents leading zeros stacking', () => {
    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('5')
    })
    expect(result.current.displayValue).toBe('5')

    act(() => {
      result.current.inputDigit('2')
    })
    expect(result.current.displayValue).toBe('52')
  })

  it('deletes digits with handleDelete', () => {
    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('4')
      result.current.inputDigit('2')
    })
    expect(result.current.displayValue).toBe('42')

    act(() => {
      result.current.handleDelete()
    })
    expect(result.current.displayValue).toBe('4')

    act(() => {
      result.current.handleDelete()
    })
    expect(result.current.displayValue).toBe('0')
  })

  it('toggles sign of active operand', () => {
    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('9')
      result.current.handleToggleSign()
    })
    expect(result.current.displayValue).toBe('-9')

    act(() => {
      result.current.handleToggleSign()
    })
    expect(result.current.displayValue).toBe('9')
  })

  it('stages binary operation and evaluates with equals', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(15)

    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('1')
      result.current.inputDigit('0')
    })

    await act(async () => {
      await result.current.stageBinary('+')
    })

    expect(result.current.subDisplayValue).toBe('10 +')

    act(() => {
      result.current.inputDigit('5')
    })

    await act(async () => {
      await result.current.executeEquals()
    })

    expect(client.calculateBinary).toHaveBeenCalledWith('+', 10, 5)
    expect(result.current.displayValue).toBe('15')
    expect(result.current.subDisplayValue).toBe('')
    expect(result.current.history.length).toBe(1)
    expect(result.current.history[0].expression).toBe('10 + 5')
    expect(result.current.history[0].result).toBe('15')
  })

  it('supports chained binary calculations', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(14)

    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('1')
      result.current.inputDigit('0')
    })

    await act(async () => {
      await result.current.stageBinary('+')
    })

    act(() => {
      result.current.inputDigit('4')
    })

    await act(async () => {
      await result.current.stageBinary('*')
    })

    expect(client.calculateBinary).toHaveBeenCalledWith('+', 10, 4)
    expect(result.current.displayValue).toBe('14')
    expect(result.current.subDisplayValue).toBe('14 *')
  })

  it('executes unary operations immediately', async () => {
    vi.mocked(client.calculateUnary).mockResolvedValueOnce(5)

    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('2')
      result.current.inputDigit('5')
    })

    await act(async () => {
      await result.current.executeUnary('sqrt')
    })

    expect(client.calculateUnary).toHaveBeenCalledWith('sqrt', 25)
    expect(result.current.displayValue).toBe('5')
    expect(result.current.history[0].expression).toBe('√(25)')
    expect(result.current.history[0].result).toBe('5')
  })

  it('handles backend 422 errors and recovers cleanly', async () => {
    vi.mocked(client.calculateBinary).mockRejectedValueOnce(
      new Error('division by zero is undefined'),
    )

    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('8')
    })

    await act(async () => {
      await result.current.stageBinary('/')
    })

    act(() => {
      result.current.inputDigit('0')
    })

    await act(async () => {
      await result.current.executeEquals()
    })

    expect(result.current.displayValue).toBe('Error')
    expect(result.current.error).toBe('division by zero is undefined')

    act(() => {
      result.current.inputDigit('3')
    })

    expect(result.current.displayValue).toBe('3')
    expect(result.current.error).toBeNull()
  })

  it('clears state on clearAll', () => {
    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.inputDigit('7')
      result.current.clearAll()
    })

    expect(result.current.displayValue).toBe('0')
    expect(result.current.subDisplayValue).toBe('')
    expect(result.current.error).toBeNull()
  })

  it('recalls history and clears history list', () => {
    const { result } = renderHook(() => useCalculator())

    act(() => {
      result.current.recallHistory({
        id: '1',
        expression: '10 * 2',
        result: '20',
        timestamp: Date.now(),
      })
    })

    expect(result.current.displayValue).toBe('20')

    act(() => {
      result.current.clearHistory()
    })

    expect(result.current.history).toEqual([])
  })
})
