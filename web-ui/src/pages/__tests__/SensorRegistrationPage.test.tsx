import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import SensorRegistrationPage from '../SensorRegistrationPage';

vi.mock('../../services/sensorsApi', () => ({
  sensorsApi: {
    listPending: vi.fn().mockResolvedValue([
      { id: 'p1', name: 'sensor-dc01', ip_address: '192.168.1.100', profile: 'datacenter_host', network_interfaces: ['eth0'], registration_key: 'REG-1', created_at: new Date().toISOString(), expires_at: new Date(Date.now()+600000).toISOString(), status: 'pending' },
    ]),
    createPending: vi.fn().mockResolvedValue({ id: 'p2' }),
    deletePending: vi.fn().mockResolvedValue(undefined),
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

describe('SensorRegistrationPage', () => {
  it('renders pending sensors', async () => {
    renderWithProviders(<SensorRegistrationPage />);
    await waitFor(() => {
      expect(screen.getByText('Sensor Registration')).toBeInTheDocument();
      expect(screen.getByText('sensor-dc01')).toBeInTheDocument();
    });
  });

  it('opens Add Sensor and creates pending', async () => {
    const { sensorsApi } = await import('../../services/sensorsApi');
    renderWithProviders(<SensorRegistrationPage />);

    await waitFor(() => expect(screen.getByText('Add Sensor')).toBeInTheDocument());
    fireEvent.click(screen.getByText('Add Sensor'));

    await waitFor(() => expect(screen.getByText('Add New Sensor')).toBeInTheDocument());

    fireEvent.change(screen.getByPlaceholderText('sensor-dc01'), { target: { value: 'sensor-cloud01' } });
    fireEvent.change(screen.getByPlaceholderText('192.168.1.100'), { target: { value: '10.0.0.1' } });

    fireEvent.click(screen.getByText('Generate Key'));

    await waitFor(() => expect(sensorsApi.createPending).toHaveBeenCalled());
  });
});
