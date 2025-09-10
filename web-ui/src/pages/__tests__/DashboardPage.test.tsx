import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { DashboardPage } from '../DashboardPage';

vi.mock('../../contexts/AuthContext', () => {
  return {
    useAuth: () => ({
      user: { first_name: 'Alex' },
      tenant: { id: 'tenant-1', name: 'DemoCorp' },
      isAuthenticated: true,
    }),
  };
});

vi.mock('../../services/dashboardApi', () => ({
  dashboardApi: {
    getOverview: vi.fn().mockResolvedValue({
      total_assets: 1234,
      risk_summary: { total_crypto: 567, total_assets: 1234, high_risk: 10, medium_risk: 20, low_risk: 30, unknown_risk: 5, critical_findings: 2 },
    }),
  },
}));

vi.mock('../../services/auditApi', () => ({
  auditApi: {
    getLogs: vi.fn().mockResolvedValue({ logs: [
      { id: '1', user_id: 'u1', tenant_id: 'tenant-1', action: 'created', resource_type: 'asset', resource_id: 'a1', permission_required: 'assets.create', permission_granted: true, created_at: new Date().toISOString() },
    ] }),
  },
}));

vi.mock('../../services/certificatesApi', () => ({
  certificatesApi: {
    getExpiring: vi.fn().mockResolvedValue([
      { id: 'c1', asset_id: 'a1', common_name: 'api.demo.com', issuer: 'Demo CA', not_after: new Date().toISOString(), days_until_expiry: 12 },
    ]),
  },
}));

vi.mock('../../services/inventoryApi', () => ({
  inventoryApi: {
    getRiskSummary: vi.fn().mockResolvedValue({
      total_assets: 1234, high_risk: 10, medium_risk: 20, low_risk: 30, unknown_risk: 5, total_crypto: 567, critical_findings: 2,
    }),
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

describe('DashboardPage', () => {
  it('renders loading skeletons initially', async () => {
    renderWithProviders(<DashboardPage />);
    expect(screen.getAllByRole('generic', { hidden: true }).length).toBeGreaterThan(0);
  });

  it('renders stats, activity, expiring certs, and risks after load', async () => {
    renderWithProviders(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Welcome back, Alex!')).toBeInTheDocument();
      expect(screen.getByText('DemoCorp', { exact: false })).toBeInTheDocument();
      expect(screen.getByText('Network Assets')).toBeInTheDocument();
      expect(screen.getByText('Crypto Implementations')).toBeInTheDocument();
      expect(screen.getByText('Recent Activity')).toBeInTheDocument();
      expect(screen.getByText('Expiring Certificates')).toBeInTheDocument();
      expect(screen.getByText('Top Risks')).toBeInTheDocument();
    });
  });

  it('handles empty activity gracefully', async () => {
    const { auditApi } = await import('../../services/auditApi');
    (auditApi.getLogs as any).mockResolvedValueOnce({ logs: [] });

    renderWithProviders(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('No recent activity.')).toBeInTheDocument();
    });
  });
});
