import api from './api';

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  domain?: string;
  subscription_tier: string;
  trial_ends_at?: string;
  billing_email?: string;
  payment_status: string;
  stripe_customer_id?: string;
  sso_enabled: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface TenantStats {
  tenant_id: string;
  tenant_name: string;
  user_count: number;
  asset_count: number;
  sensor_count: number;
  last_activity?: string;
  storage_used: number;
  api_requests: number;
}

export interface TenantsResponse {
  tenants: Tenant[];
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
}

export const tenantsApi = {
  getTenants: async (page = 1, limit = 20): Promise<TenantsResponse> => {
    const response = await api.get(`/admin/tenants?page=${page}&limit=${limit}`);
    return response.data;
  },

  getTenant: async (id: string): Promise<{ tenant: Tenant }> => {
    const response = await api.get(`/admin/tenants/${id}`);
    return response.data;
  },

  createTenant: async (tenant: Partial<Tenant>): Promise<{ message: string; tenant_id: string }> => {
    const response = await api.post('/admin/tenants', tenant);
    return response.data;
  },

  updateTenant: async (id: string, tenant: Partial<Tenant>): Promise<{ message: string }> => {
    const response = await api.put(`/admin/tenants/${id}`, tenant);
    return response.data;
  },

  deleteTenant: async (id: string): Promise<{ message: string }> => {
    const response = await api.delete(`/admin/tenants/${id}`);
    return response.data;
  },

  suspendTenant: async (id: string): Promise<{ message: string }> => {
    const response = await api.post(`/admin/tenants/${id}/suspend`);
    return response.data;
  },

  activateTenant: async (id: string): Promise<{ message: string }> => {
    const response = await api.post(`/admin/tenants/${id}/activate`);
    return response.data;
  },

  getTenantStats: async (id: string): Promise<{ stats: TenantStats }> => {
    const response = await api.get(`/admin/tenants/${id}/stats`);
    return response.data;
  },

  getTenantsStats: async (): Promise<{ tenants_stats: TenantStats[] }> => {
    const response = await api.get('/admin/stats/tenants');
    return response.data;
  },
};
