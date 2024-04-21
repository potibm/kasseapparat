import { render, screen } from '@testing-library/react'
import App from './App'

test('renders the application', () => {
  render(<App />)
  const linkElement = screen.getByText(/Checkout/i)
  expect(linkElement).toBeInTheDocument()
})
