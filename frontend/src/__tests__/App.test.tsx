import { render } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { BrowserRouter } from 'react-router-dom'
import App from '@/App'

// Wrapper for tests that need routing
const AppWithRouter = () => (
  <BrowserRouter>
    <App />
  </BrowserRouter>
)

describe('App', () => {
  it('renders without crashing', () => {
    const { container } = render(<AppWithRouter />)
    // App should render without throwing errors
    expect(container).toBeInTheDocument()
  })

  it('should pass basic test', () => {
    expect(1 + 1).toBe(2)
  })
})
