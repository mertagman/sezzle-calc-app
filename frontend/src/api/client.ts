import type {
  BinaryOperator,
  BinaryRequest,
  CalculateResponse,
  ErrorResponse,
  UnaryOperator,
  UnaryRequest,
} from './types'

const BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

const BINARY_ENDPOINTS: Record<BinaryOperator, string> = {
  '+': `${BASE_URL}/api/v1/add`,
  '-': `${BASE_URL}/api/v1/subtract`,
  '*': `${BASE_URL}/api/v1/multiply`,
  '/': `${BASE_URL}/api/v1/divide`,
  '^': `${BASE_URL}/api/v1/power`,
}

const UNARY_ENDPOINTS: Record<UnaryOperator, string> = {
  sqrt: `${BASE_URL}/api/v1/sqrt`,
  '%': `${BASE_URL}/api/v1/percentage`,
}

export async function calculateBinary(
  op: BinaryOperator,
  a: number,
  b: number,
): Promise<number> {
  const endpoint = BINARY_ENDPOINTS[op]
  if (!endpoint) {
    throw new Error(`Unsupported binary operator: ${op}`)
  }

  const payload: BinaryRequest = { a, b }

  const response = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    const errorData = (await response.json().catch(() => null)) as ErrorResponse | null
    throw new Error(errorData?.error || `Calculation failed (${response.status})`)
  }

  const data = (await response.json()) as CalculateResponse
  return data.result
}

export async function calculateUnary(
  op: UnaryOperator,
  a: number,
): Promise<number> {
  const endpoint = UNARY_ENDPOINTS[op]
  if (!endpoint) {
    throw new Error(`Unsupported unary operator: ${op}`)
  }

  const payload: UnaryRequest = { a }

  const response = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    const errorData = (await response.json().catch(() => null)) as ErrorResponse | null
    throw new Error(errorData?.error || `Calculation failed (${response.status})`)
  }

  const data = (await response.json()) as CalculateResponse
  return data.result
}

export async function checkHealth(): Promise<boolean> {
  try {
    const response = await fetch(`${BASE_URL}/health`)
    return response.ok
  } catch {
    return false
  }
}
