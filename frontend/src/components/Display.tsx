interface DisplayProps {
  displayValue: string
  subDisplayValue: string
  error: string | null
  isCalculating: boolean
}

export function Display({
  displayValue,
  subDisplayValue,
  error,
  isCalculating,
}: DisplayProps) {
  const isLong = displayValue.length > 11
  const isVeryLong = displayValue.length > 14

  return (
    <div className="retro-display" role="region" aria-label="Calculator display">
      <div className="sub-display" aria-live="polite">
        {subDisplayValue ? subDisplayValue : '\u00A0'}
        {isCalculating && <span className="calculating-indicator">...</span>}
      </div>
      <div
        className={`main-display ${isVeryLong ? 'text-xs' : isLong ? 'text-sm' : ''} ${error ? 'has-error' : ''}`}
        aria-live="assertive"
        data-testid="main-display"
      >
        {error ? error : displayValue}
      </div>
    </div>
  )
}
