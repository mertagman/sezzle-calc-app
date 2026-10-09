import type { BinaryOperator, UnaryOperator } from '../api/types'

interface KeypadProps {
  onDigit: (digit: string) => void
  onDecimal: () => void
  onDelete: () => void
  onToggleSign: () => void
  onClear: () => void
  onBinary: (op: BinaryOperator) => void
  onUnary: (op: UnaryOperator) => void
  onEquals: () => void
  disabled?: boolean
}

export function Keypad({
  onDigit,
  onDecimal,
  onDelete,
  onToggleSign,
  onClear,
  onBinary,
  onUnary,
  onEquals,
  disabled = false,
}: KeypadProps) {
  return (
    <div className="retro-keypad" role="group" aria-label="Calculator keypad">
      <button
        type="button"
        className="retro-btn btn-clear"
        onClick={onClear}
        data-key="c"
        aria-label="Clear all"
      >
        C
      </button>
      <button
        type="button"
        className="retro-btn btn-del"
        onClick={onDelete}
        data-key="del"
        aria-label="Delete last digit"
      >
        DEL
      </button>
      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onUnary('sqrt')}
        disabled={disabled}
        data-key="sqrt"
        aria-label="Square root"
      >
        √
      </button>
      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onUnary('%')}
        disabled={disabled}
        data-key="percent"
        aria-label="Percentage"
      >
        %
      </button>

      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onBinary('^')}
        disabled={disabled}
        data-key="pow"
        aria-label="Exponentiation"
      >
        xʸ
      </button>
      <button
        type="button"
        className="retro-btn btn-dummy"
        disabled
        aria-hidden="true"
        tabIndex={-1}
      >
        —
      </button>
      <button
        type="button"
        className="retro-btn btn-dummy"
        disabled
        aria-hidden="true"
        tabIndex={-1}
      >
        —
      </button>
      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onBinary('/')}
        disabled={disabled}
        data-key="/"
        aria-label="Divide"
      >
        ÷
      </button>

      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('7')}
        data-key="7"
      >
        7
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('8')}
        data-key="8"
      >
        8
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('9')}
        data-key="9"
      >
        9
      </button>
      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onBinary('*')}
        disabled={disabled}
        data-key="*"
        aria-label="Multiply"
      >
        ×
      </button>

      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('4')}
        data-key="4"
      >
        4
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('5')}
        data-key="5"
      >
        5
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('6')}
        data-key="6"
      >
        6
      </button>
      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onBinary('-')}
        disabled={disabled}
        data-key="-"
        aria-label="Subtract"
      >
        −
      </button>

      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('1')}
        data-key="1"
      >
        1
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('2')}
        data-key="2"
      >
        2
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('3')}
        data-key="3"
      >
        3
      </button>
      <button
        type="button"
        className="retro-btn btn-op"
        onClick={() => onBinary('+')}
        disabled={disabled}
        data-key="+"
        aria-label="Add"
      >
        +
      </button>

      <button
        type="button"
        className="retro-btn btn-num"
        onClick={() => onDigit('0')}
        data-key="0"
      >
        0
      </button>
      <button
        type="button"
        className="retro-btn btn-num"
        onClick={onDecimal}
        data-key="."
        aria-label="Decimal point"
      >
        .
      </button>
      <button
        type="button"
        className="retro-btn btn-fn"
        onClick={onToggleSign}
        data-key="+/-"
        aria-label="Negate number"
      >
        ±
      </button>
      <button
        type="button"
        className="retro-btn btn-equals"
        onClick={onEquals}
        disabled={disabled}
        data-key="="
        aria-label="Equals"
      >
        =
      </button>
    </div>
  )
}
