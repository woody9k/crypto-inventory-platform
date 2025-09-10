import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import ReportsPage from '../ReportsPage';

vi.mock('../../services/reportsApi', () => ({
  reportsApi: {
    list: vi.fn().mockResolvedValue([
      { id: 'r1', title: 'Crypto Summary Report', type: 'crypto_summary', status: 'completed', created_at: new Date().toISOString(), completed_at: new Date().toISOString(), download_url: '/api/v1/reports/r1/download' },
    ]),
    templates: vi.fn().mockResolvedValue([
      { id: 'crypto_summary', name: 'Crypto Summary Report', description: 'Overview', type: 'summary', category: 'crypto' },
    ]),
    generate: vi.fn().mockResolvedValue({ id: 'r2' }),
    delete: vi.fn().mockResolvedValue(undefined),
  },
}));

const renderWithProviders = (ui: React.ReactElement) => {
  const client = new QueryClient();
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        {ui}
      </MemoryRouter>
    </QueryClientProvider>
  );
};

describe('ReportsPage', () => {
  it('renders reports list and stats', async () => {
    renderWithProviders(<ReportsPage />);
    await waitFor(() => {
      expect(screen.getByText('Reports')).toBeInTheDocument();
      expect(screen.getByText('Total Reports')).toBeInTheDocument();
      expect(screen.getByText('Crypto Summary Report')).toBeInTheDocument();
    });
  });

  it('opens generate modal and triggers generate', async () => {
    const { reportsApi } = await import('../../services/reportsApi');
    renderWithProviders(<ReportsPage />);

    await waitFor(() => expect(screen.getByText('Generate Report')).toBeInTheDocument());
    fireEvent.click(screen.getByText('Generate Report'));

    await waitFor(() => expect(screen.getByText('Generate New Report')).toBeInTheDocument());

    // Click the first template
    fireEvent.click(screen.getByText('Crypto Summary Report'));

    await waitFor(() => expect(reportsApi.generate).toHaveBeenCalled());
  });
});
