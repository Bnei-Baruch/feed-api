import React from 'react';
import { render } from '@testing-library/react';
import App from './App';

test('renders feed viewer', () => {
  const { getByText } = render(<App />);
  const pElement = getByText(/Feed Viewer/i);
  expect(pElement).toBeInTheDocument();
});
