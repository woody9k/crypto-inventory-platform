import api from './api';
import { Asset, AssetsResponse, RiskSummary, AssetFilters, CryptoImplementation } from '../types/inventory';

const INVENTORY_BASE = import.meta.env.VITE_INVENTORY_URL || 'http://localhost:8082';

export const inventoryApi = {
  // Get assets with filtering and pagination
  getAssets: async (filters?: AssetFilters): Promise<AssetsResponse> => {
    const params = new URLSearchParams();
    
    if (filters?.search) params.append('search', filters.search);
    if (filters?.asset_type?.length) filters.asset_type.forEach(type => params.append('asset_type', type));
    if (filters?.environment?.length) filters.environment.forEach(env => params.append('environment', env));
    if (filters?.risk_level?.length) filters.risk_level.forEach(level => params.append('risk_level', level));
    if (filters?.protocol?.length) filters.protocol.forEach(protocol => params.append('protocol', protocol));
    if (filters?.business_unit?.length) filters.business_unit.forEach(unit => params.append('business_unit', unit));
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.page_size) params.append('page_size', filters.page_size.toString());
    if (filters?.sort_by) params.append('sort_by', filters.sort_by);
    if (filters?.sort_order) params.append('sort_order', filters.sort_order);

    const response = await api.get(`${INVENTORY_BASE}/api/v1/assets?${params.toString()}`);
    return response.data;
  },

  // Get single asset by ID
  getAsset: async (assetId: string): Promise<Asset> => {
    const response = await api.get(`${INVENTORY_BASE}/api/v1/assets/${assetId}`);
    return response.data.asset;
  },

  // Get crypto implementations for an asset
  getAssetCrypto: async (assetId: string): Promise<CryptoImplementation[]> => {
    const response = await api.get(`${INVENTORY_BASE}/api/v1/assets/${assetId}/crypto`);
    return response.data.crypto_implementations;
  },

  // Search assets
  searchAssets: async (query: string, limit?: number): Promise<{ assets: Asset[]; total: number; showing: number }> => {
    const params = new URLSearchParams();
    params.append('q', query);
    if (limit) params.append('limit', limit.toString());

    const response = await api.get(`${INVENTORY_BASE}/api/v1/assets/search?${params.toString()}`);
    return response.data;
  },

  // Get risk summary
  getRiskSummary: async (): Promise<RiskSummary> => {
    const response = await api.get(`${INVENTORY_BASE}/api/v1/risk/summary`);
    return response.data.risk_summary;
  },
};
