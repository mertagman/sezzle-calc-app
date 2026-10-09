export type BinaryOperator = '+' | '-' | '*' | '/' | '^'
export type UnaryOperator = 'sqrt' | '%'

export interface BinaryRequest {
  a: number
  b: number
}

export interface UnaryRequest {
  a: number
}

export interface CalculateResponse {
  result: number
}

export interface ErrorResponse {
  error: string
}

export interface HistoryItem {
  id: string
  expression: string
  result: string
  timestamp: number
}
