import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import * as client from '../api/client'
import { Calculator } from './Calculator'

vi.mock('../api/client', () => ({
  calculateBinary: vi.fn(),
  calculateUnary: vi.fn(),
  checkHealth: vi.fn(),
}))

describe('Calculator Component Integration', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('renders display, buttons, and history panel', () => {
    render(<Calculator />)

    expect(screen.getByTestId('main-display')).toHaveTextContent('0')
    expect(screen.getByRole('button', { name: 'Clear all' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Delete last digit' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'HISTORY' })).toBeInTheDocument()
  })

  it('inputs digits and deletes with DEL button', () => {
    render(<Calculator />)

    fireEvent.click(screen.getByText('7'))
    fireEvent.click(screen.getByText('8'))
    expect(screen.getByTestId('main-display')).toHaveTextContent('78')

    fireEvent.click(screen.getByRole('button', { name: 'Delete last digit' }))
    expect(screen.getByTestId('main-display')).toHaveTextContent('7')
  })

  it('executes addition and adds item to history', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(42)

    render(<Calculator />)

    fireEvent.click(screen.getByText('4'))
    fireEvent.click(screen.getByText('0'))
    fireEvent.click(screen.getByRole('button', { name: 'Add' }))
    fireEvent.click(screen.getByText('2'))
    fireEvent.click(screen.getByRole('button', { name: 'Equals' }))

    await waitFor(() => {
      expect(screen.getByTestId('main-display')).toHaveTextContent('42')
    })

    expect(screen.getByText('40 + 2 =')).toBeInTheDocument()
    expect(screen.getByText('42')).toBeInTheDocument()
  })

  it('recalls result into display when history item is clicked', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(99)

    render(<Calculator />)

    fireEvent.click(screen.getByText('9'))
    fireEvent.click(screen.getByRole('button', { name: 'Add' }))
    fireEvent.click(screen.getByText('9'))
    fireEvent.click(screen.getByRole('button', { name: 'Equals' }))

    await waitFor(() => {
      expect(screen.getByText('9 + 9 =')).toBeInTheDocument()
    })

    fireEvent.click(screen.getByRole('button', { name: 'Clear all' }))
    expect(screen.getByTestId('main-display')).toHaveTextContent('0')

    fireEvent.click(screen.getByRole('button', { name: 'Click to recall result' }))
    expect(screen.getByTestId('main-display')).toHaveTextContent('99')
  })

  it('handles physical keyboard events', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(8)

    render(<Calculator />)

    fireEvent.keyDown(window, { key: '5' })
    fireEvent.keyDown(window, { key: '+' })
    fireEvent.keyDown(window, { key: '3' })
    fireEvent.keyDown(window, { key: 'Enter' })

    await waitFor(() => {
      expect(screen.getByTestId('main-display')).toHaveTextContent('8')
    })

    fireEvent.keyDown(window, { key: 'Backspace' })
    expect(screen.getByTestId('main-display')).toHaveTextContent('8')
  })
})
