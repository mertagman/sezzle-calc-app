import { useCallback, useEffect, useState } from 'react'
import { calculateBinary, calculateUnary } from '../api/client'
import type { BinaryOperator, HistoryItem, UnaryOperator } from '../api/types'
import {
  appendDecimal,
  appendDigit,
  applyBackspace,
  formatNumberResult,
  toggleSign,
} from '../utils/format'

const STORAGE_KEY = 'retro_calc_history'

function loadHistoryFromStorage(): HistoryItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

export function useCalculator() {
  const [displayValue, setDisplayValue] = useState<string>('0')
  const [subDisplayValue, setSubDisplayValue] = useState<string>('')
  const [accumulator, setAccumulator] = useState<number | null>(null)
  const [pendingOperator, setPendingOperator] = useState<BinaryOperator | null>(null)
  const [waitingForOperand, setWaitingForOperand] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)
  const [isCalculating, setIsCalculating] = useState<boolean>(false)
  const [history, setHistory] = useState<HistoryItem[]>(loadHistoryFromStorage)

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(history))
    } catch {
      // Storage unavailable or full
    }
  }, [history])

  const appendHistoryItem = useCallback((expression: string, result: string) => {
    const newItem: HistoryItem = {
      id: `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`,
      expression,
      result,
      timestamp: Date.now(),
    }
    setHistory((prev) => [newItem, ...prev].slice(0, 50))
  }, [])

  const inputDigit = useCallback((digit: string) => {
    if (error) {
      setError(null)
      setDisplayValue(digit)
      setWaitingForOperand(false)
      return
    }

    if (waitingForOperand) {
      setDisplayValue(digit)
      setWaitingForOperand(false)
    } else {
      setDisplayValue((prev) => appendDigit(prev, digit))
    }
  }, [error, waitingForOperand])

  const inputDecimal = useCallback(() => {
    if (error) {
      setError(null)
      setDisplayValue('0.')
      setWaitingForOperand(false)
      return
    }

    if (waitingForOperand) {
      setDisplayValue('0.')
      setWaitingForOperand(false)
    } else {
      setDisplayValue((prev) => appendDecimal(prev))
    }
  }, [error, waitingForOperand])

  const handleDelete = useCallback(() => {
    if (error) {
      setError(null)
      setDisplayValue('0')
      return
    }

    if (waitingForOperand) {
      return
    }

    setDisplayValue((prev) => applyBackspace(prev))
  }, [error, waitingForOperand])

  const handleToggleSign = useCallback(() => {
    if (error) return
    setDisplayValue((prev) => toggleSign(prev))
  }, [error])

  const executeUnary = useCallback(async (op: UnaryOperator) => {
    if (error) return

    const currentNum = parseFloat(displayValue)
    if (Number.isNaN(currentNum)) return

    setIsCalculating(true)
    try {
      const res = await calculateUnary(op, currentNum)
      const formatted = formatNumberResult(res)
      setDisplayValue(formatted)
      setWaitingForOperand(true)

      const expr = op === 'sqrt' ? `√(${displayValue})` : `${displayValue}%`
      appendHistoryItem(expr, formatted)
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Error'
      setError(message)
      setDisplayValue('Error')
    } finally {
      setIsCalculating(false)
    }
  }, [appendHistoryItem, displayValue, error])

  const stageBinary = useCallback(async (op: BinaryOperator) => {
    if (error) return

    const currentNum = parseFloat(displayValue)
    if (Number.isNaN(currentNum)) return

    if (pendingOperator !== null && !waitingForOperand && accumulator !== null) {
      setIsCalculating(true)
      try {
        const res = await calculateBinary(pendingOperator, accumulator, currentNum)
        const formatted = formatNumberResult(res)
        setAccumulator(res)
        setDisplayValue(formatted)
        setSubDisplayValue(`${formatted} ${op}`)
        setPendingOperator(op)
        setWaitingForOperand(true)
      } catch (err: unknown) {
        const message = err instanceof Error ? err.message : 'Error'
        setError(message)
        setDisplayValue('Error')
        setSubDisplayValue('')
        setAccumulator(null)
        setPendingOperator(null)
      } finally {
        setIsCalculating(false)
      }
      return
    }

    setAccumulator(currentNum)
    setPendingOperator(op)
    setSubDisplayValue(`${displayValue} ${op}`)
    setWaitingForOperand(true)
  }, [accumulator, displayValue, error, pendingOperator, waitingForOperand])

  const executeEquals = useCallback(async () => {
    if (error) return
    if (pendingOperator === null || accumulator === null) return

    const currentNum = parseFloat(displayValue)
    if (Number.isNaN(currentNum)) return

    const activeOp = pendingOperator
    const activeAcc = accumulator

    setIsCalculating(true)
    try {
      const res = await calculateBinary(activeOp, activeAcc, currentNum)
      const formatted = formatNumberResult(res)
      setDisplayValue(formatted)
      setSubDisplayValue('')
      setAccumulator(null)
      setPendingOperator(null)
      setWaitingForOperand(true)

      appendHistoryItem(`${activeAcc} ${activeOp} ${currentNum}`, formatted)
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Error'
      setError(message)
      setDisplayValue('Error')
      setSubDisplayValue('')
      setAccumulator(null)
      setPendingOperator(null)
    } finally {
      setIsCalculating(false)
    }
  }, [accumulator, appendHistoryItem, displayValue, error, pendingOperator])

  const clearAll = useCallback(() => {
    setDisplayValue('0')
    setSubDisplayValue('')
    setAccumulator(null)
    setPendingOperator(null)
    setWaitingForOperand(false)
    setError(null)
    setIsCalculating(false)
  }, [])

  const clearHistory = useCallback(() => {
    setHistory([])
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch {
      // Storage unavailable
    }
  }, [])

  const recallHistory = useCallback((item: HistoryItem) => {
    setDisplayValue(item.result)
    setWaitingForOperand(true)
    setError(null)
  }, [])

  return {
    displayValue,
    subDisplayValue,
    accumulator,
    pendingOperator,
    waitingForOperand,
    error,
    isCalculating,
    history,
    inputDigit,
    inputDecimal,
    handleDelete,
    handleToggleSign,
    executeUnary,
    stageBinary,
    executeEquals,
    clearAll,
    clearHistory,
    recallHistory,
  }
}
