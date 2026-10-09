import type { HistoryItem } from '../api/types'

interface HistoryProps {
  items: HistoryItem[]
  onSelect: (item: HistoryItem) => void
  onClear: () => void
}

export function History({ items, onSelect, onClear }: HistoryProps) {
  return (
    <aside className="retro-history" aria-label="Calculation history">
      <div className="history-header">
        <h2 className="history-title">HISTORY</h2>
        <button
          type="button"
          className="history-clear-btn"
          onClick={onClear}
          disabled={items.length === 0}
          aria-label="Clear History"
        >
          Clear History
        </button>
      </div>

      <div className="history-tape">
        {items.length === 0 ? (
          <div className="history-empty">No calculations yet</div>
        ) : (
          <ul className="history-list">
            {items.map((item) => (
              <li key={item.id} className="history-item">
                <button
                  type="button"
                  className="history-item-btn"
                  onClick={() => onSelect(item)}
                  title="Click to recall result"
                >
                  <span className="history-expr">{item.expression} =</span>
                  <span className="history-res">{item.result}</span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </aside>
  )
}
