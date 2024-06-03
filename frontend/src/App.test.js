import React from 'react'
import { render, screen } from '@testing-library/react'
import App from './App'
import { test, expect } from '@jest/globals'

test('renders the application', () => {
  render(<App />)
  const headlineElement = screen.getByText(/Kasseapparat/i)
  expect(headlineElement).toBeInTheDocument()
})
