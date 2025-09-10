// Audit API client: fetches tenant-scoped permission/activity logs for dashboard Recent Activity.
import api from './api';

export interface AuditLog {
  id: string;
  user_id: string;
  tenant_id: string;
  action: string;
  resource_type: string;
  resource_id: string;
  permission_required: string;
  permission_granted: boolean;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface GetAuditLogsParams {
  page?: number;
  limit?: number;
  user_id?: string;
  tenant_id?: string;
  from?: string; // ISO datetime
  to?: string;   // ISO datetime
}

export interface GetAuditLogsResponse {
  logs: AuditLog[];
}

export const auditApi = {
  async getLogs(params: GetAuditLogsParams = {}): Promise<GetAuditLogsResponse> {
    const search = new URLSearchParams();
    if (params.page) search.append('page', String(params.page));
    if (params.limit) search.append('limit', String(params.limit));
    if (params.user_id) search.append('user_id', params.user_id);
    if (params.tenant_id) search.append('tenant_id', params.tenant_id);
    if (params.from) search.append('from', params.from);
    if (params.to) search.append('to', params.to);

    const response = await api.get(`/audit/logs?${search.toString()}`);
    const logs = Array.isArray(response.data?.logs) ? response.data.logs : (response.data || []);
    return { logs };
  },
};
