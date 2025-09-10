// Sensors API client: integrates with sensor-manager service (default http://localhost:8085/api/v1)
// Supports pending registration key lifecycle and sensor registration.
import api from './api';

const SENSOR_BASE = 'http://localhost:8085/api/v1';

export interface PendingSensorPayload {
  name: string;
  ip_address: string;
  profile: string;
  network_interfaces?: string[];
  tags?: string[];
}

export interface PendingSensorItem {
  id: string;
  name: string;
  ip_address: string;
  profile: string;
  network_interfaces: string[];
  registration_key: string;
  created_at: string;
  expires_at: string;
  status: 'pending' | 'expired' | 'used';
}

export interface SensorItem {
  id: string;
  name: string;
  status: 'active' | 'inactive' | 'error' | 'unknown' | 'pending';
  platform?: string;
  version?: string;
  profile?: string;
  last_seen_at?: string;
  ip_address?: string;
}

export const sensorsApi = {
  async createPending(payload: PendingSensorPayload): Promise<PendingSensorItem> {
    const res = await api.post(`${SENSOR_BASE}/sensors/pending`, payload);
    return res.data.pending_sensor ?? res.data;
    },
  async listPending(): Promise<PendingSensorItem[]> {
    const res = await api.get(`${SENSOR_BASE}/sensors/pending`);
    return res.data.pending_sensors ?? res.data ?? [];
  },
  async deletePending(key: string): Promise<void> {
    await api.delete(`${SENSOR_BASE}/sensors/pending/${encodeURIComponent(key)}`);
  },
  async registerSensor(payload: { registration_key: string }): Promise<{ sensor: SensorItem }> {
    const res = await api.post(`${SENSOR_BASE}/sensors/register`, payload);
    return res.data;
  },
};
