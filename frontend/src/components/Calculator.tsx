import { useEffect } from 'react'
import { useCalculator } from '../hooks/useCalculator'
import { Display } from './Display'
import { History } from './History'
import { Keypad } from './Keypad'

export function Calculator() {
  const {
    displayValue,
    subDisplayValue,
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
  } = useCalculator()

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
        return
      }

      const key = e.key

      if (/^[0-9]$/.test(key)) {
        e.preventDefault()
        inputDigit(key)
        return
      }

      if (key === '.') {
        e.preventDefault()
        inputDecimal()
        return
      }

      if (key === 'Backspace' || key === 'Delete') {
        e.preventDefault()
        handleDelete()
        return
      }

      if (key === 'Escape' || key === 'c' || key === 'C') {
        e.preventDefault()
        clearAll()
        return
      }

      if (key === '+' || key === '-' || key === '*' || key === '/' || key === '^') {
        e.preventDefault()
        stageBinary(key)
        return
      }

      if (key === '%') {
        e.preventDefault()
        executeUnary('%')
        return
      }

      if (key === '=' || key === 'Enter') {
        e.preventDefault()
        executeEquals()
        return
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [
    clearAll,
    executeEquals,
    executeUnary,
    handleDelete,
    inputDecimal,
    inputDigit,
    stageBinary,
  ])

  return (
    <div className="retro-calc-layout">
      <main className="retro-chassis">
        <header className="chassis-header">
          <h1 className="brand-badge">CALCULATOR</h1>
        </header>

        <Display
          displayValue={displayValue}
          subDisplayValue={subDisplayValue}
          error={error}
          isCalculating={isCalculating}
        />

        <Keypad
          onDigit={inputDigit}
          onDecimal={inputDecimal}
          onDelete={handleDelete}
          onToggleSign={handleToggleSign}
          onClear={clearAll}
          onBinary={stageBinary}
          onUnary={executeUnary}
          onEquals={executeEquals}
          disabled={isCalculating}
        />
      </main>

      <History
        items={history}
        onSelect={recallHistory}
        onClear={clearHistory}
      />
    </div>
  )
}
