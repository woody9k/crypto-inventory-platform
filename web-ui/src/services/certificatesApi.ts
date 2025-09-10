import api from './api';

export interface ExpiringCertificate {
  id: string;
  asset_id: string;
  common_name?: string;
  issuer?: string;
  not_after: string; // expiry
  days_until_expiry: number;
}

export const certificatesApi = {
  async getExpiring(limit: number = 5): Promise<ExpiringCertificate[]> {
    // Assuming inventory service exposes this; adjust path if needed
    const response = await api.get(`http://localhost:8082/api/v1/certificates/expiring?limit=${limit}`);
    const list = Array.isArray(response.data?.certificates) ? response.data.certificates : (response.data || []);
    return list;
  },
};
