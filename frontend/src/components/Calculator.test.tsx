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
    vi.resetAllMocks()
  })

  it('renders display, buttons, and history panel', () => {
    render(<Calculator />)

    expect(screen.getByTestId('main-display')).toHaveTextContent('0')
    expect(screen.getByRole('button', { name: 'Clear all' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Delete last digit' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'HISTORY' })).toBeInTheDocument()
  })

  it('inputs digits and deletes with DEL button', () => {
    const { container } = render(<Calculator />)

    const btn7 = container.querySelector('[data-key="7"]') as HTMLButtonElement
    const btn8 = container.querySelector('[data-key="8"]') as HTMLButtonElement
    const btnDel = screen.getByRole('button', { name: 'Delete last digit' })

    fireEvent.click(btn7)
    fireEvent.click(btn8)
    expect(screen.getByTestId('main-display')).toHaveTextContent('78')

    fireEvent.click(btnDel)
    expect(screen.getByTestId('main-display')).toHaveTextContent('7')
  })

  it('executes addition and adds item to history', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(42)

    const { container } = render(<Calculator />)

    const btn4 = container.querySelector('[data-key="4"]') as HTMLButtonElement
    const btn0 = container.querySelector('[data-key="0"]') as HTMLButtonElement
    const btn2 = container.querySelector('[data-key="2"]') as HTMLButtonElement
    const btnAdd = screen.getByRole('button', { name: 'Add' })
    const btnEquals = screen.getByRole('button', { name: 'Equals' })

    fireEvent.click(btn4)
    fireEvent.click(btn0)
    fireEvent.click(btnAdd)
    fireEvent.click(btn2)
    fireEvent.click(btnEquals)

    await waitFor(() => {
      expect(screen.getByTestId('main-display')).toHaveTextContent('42')
    })

    expect(screen.getByText('40 + 2 =')).toBeInTheDocument()
  })

  it('recalls result into display when history item is clicked', async () => {
    vi.mocked(client.calculateBinary).mockResolvedValueOnce(99)

    const { container } = render(<Calculator />)

    const btn9 = container.querySelector('[data-key="9"]') as HTMLButtonElement
    const btnAdd = screen.getByRole('button', { name: 'Add' })
    const btnEquals = screen.getByRole('button', { name: 'Equals' })
    const btnClear = screen.getByRole('button', { name: 'Clear all' })

    fireEvent.click(btn9)
    fireEvent.click(btnAdd)
    fireEvent.click(btn9)
    fireEvent.click(btnEquals)

    await waitFor(() => {
      expect(screen.getByText('9 + 9 =')).toBeInTheDocument()
    })

    fireEvent.click(btnClear)
    expect(screen.getByTestId('main-display')).toHaveTextContent('0')

    const recallBtn = screen.getByRole('button', { name: /9 \+ 9/ })
    fireEvent.click(recallBtn)
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
