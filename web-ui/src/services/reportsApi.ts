// Reports API client: integrates with report-generator service via API gateway
// Routes through API gateway at http://localhost:8080/api/v1/reports/ for proper load balancing
// and authentication. Provides list, get, generate, delete, and templates operations used by ReportsPage.
import api from './api';

export interface ReportItem {
  id: string;
  title: string;
  type: string;
  status: 'generating' | 'completed' | 'failed';
  created_at: string;
  completed_at?: string;
  download_url?: string;
  data?: any;
}

export interface ReportTemplateItem {
  id: string;
  name: string;
  description: string;
  type: string;
  category: string;
}

// API gateway endpoint for reports service
const REPORTS_BASE = 'http://localhost:8080/api/v1';

export const reportsApi = {
  async list(): Promise<ReportItem[]> {
    const res = await api.get(`${REPORTS_BASE}/reports`);
    return res.data.reports ?? res.data ?? [];
  },
  async get(id: string): Promise<ReportItem> {
    const res = await api.get(`${REPORTS_BASE}/reports/${id}`);
    return res.data;
  },
  async delete(id: string): Promise<void> {
    await api.delete(`${REPORTS_BASE}/reports/${id}`);
  },
  async generate(payload: { type: string; title: string; parameters?: any; format?: string }): Promise<{ id: string }> {
    const res = await api.post(`${REPORTS_BASE}/reports/generate`, payload);
    return res.data;
  },
  async templates(): Promise<ReportTemplateItem[]> {
    const res = await api.get(`${REPORTS_BASE}/reports/templates`);
    return res.data.templates ?? res.data ?? [];
  },
  async download(id: string, format: 'pdf' | 'excel' | 'json'): Promise<Blob> {
    const res = await api.get(`${REPORTS_BASE}/reports/${id}/download?format=${format}` as any, {
      responseType: 'blob',
    } as any);
    return res.data as Blob;
  },
};
