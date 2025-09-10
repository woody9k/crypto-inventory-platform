import api from './api';

export interface Role {
  id: string;
  name: string;
  display_name: string;
  description: string;
  is_system_role?: boolean;
}

export interface PermissionItem {
  id: string;
  name: string;
  resource: string;
  action: string;
  scope?: string;
  description?: string;
}

export const rbacApi = {
  async getTenantRoles(tenantId: string): Promise<Role[]> {
    const res = await api.get(`/tenant/${tenantId}/roles`);
    return res.data.roles ?? res.data ?? [];
  },

  async getTenantPermissions(): Promise<PermissionItem[]> {
    const res = await api.get('/permissions');
    return res.data.permissions ?? res.data ?? [];
  },

  async getPermissionMatrix(tenantId: string, roleId: string): Promise<any> {
    const res = await api.get(`/tenant/${tenantId}/roles/${roleId}/matrix`);
    return res.data;
  },

  async updateRolePermissions(tenantId: string, roleId: string, permissionIds: string[]): Promise<void> {
    await api.put(`/tenant/${tenantId}/roles/${roleId}/permissions`, { permission_ids: permissionIds });
  },

  async getUserRoles(tenantId: string, userId: string): Promise<Role[]> {
    const res = await api.get(`/tenant/${tenantId}/users/${userId}/roles`);
    return res.data.roles ?? res.data ?? [];
  },

  async assignUserRole(tenantId: string, userId: string, roleId: string): Promise<void> {
    await api.post(`/tenant/${tenantId}/users/${userId}/roles`, { role_id: roleId });
  },

  async removeUserRole(tenantId: string, userId: string, roleId: string): Promise<void> {
    await api.delete(`/tenant/${tenantId}/users/${userId}/roles/${roleId}`);
  },
};
