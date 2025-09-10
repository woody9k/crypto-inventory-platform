import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import RoleManagement from '../RoleManagement';

vi.mock('../../../contexts/AuthContext', () => ({
  useAuth: () => ({ tenant: { id: 'tenant-1' } }),
}));

vi.mock('../../../services/rbacApi', () => ({
  rbacApi: {
    getTenantRoles: vi.fn().mockResolvedValue([
      { id: 'role-1', name: 'viewer', display_name: 'Viewer', description: 'Read-only', is_system_role: true },
    ]),
    getTenantPermissions: vi.fn().mockResolvedValue([
      { id: 'perm-1', name: 'assets.read', resource: 'assets', action: 'read', description: 'View assets' },
    ]),
    updateRolePermissions: vi.fn().mockResolvedValue(undefined),
  },
}));

const renderWithProviders = (ui: React.ReactElement) => {
  const client = new QueryClient();
  return render(
    <QueryClientProvider client={client}>
      {ui}
    </QueryClientProvider>
  );
};

describe('RoleManagement', () => {
  it('loads roles and permissions and toggles permission', async () => {
    const { rbacApi } = await import('../../../services/rbacApi');
    renderWithProviders(<RoleManagement />);

    await waitFor(() => {
      expect(screen.getByText('Role-Based Access Control')).toBeInTheDocument();
      expect(screen.getByText('Viewer')).toBeInTheDocument();
    });

    // Select role
    fireEvent.click(screen.getByText('Viewer'));

    await waitFor(() => expect(screen.getByText('Permissions for Viewer')).toBeInTheDocument());

    // Toggle permission
    fireEvent.click(screen.getByText('Denied'));

    await waitFor(() => expect(rbacApi.updateRolePermissions).toHaveBeenCalled());
  });
});
