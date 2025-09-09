export interface Asset {
  id: string;
  tenant_id: string;
  hostname?: string;
  ip_address?: string;
  port?: number;
  asset_type: string;
  operating_system?: string;
  environment?: string;
  business_unit?: string;
  owner_email?: string;
  description?: string;
  tags: Record<string, any>;
  metadata: Record<string, any>;
  first_discovered_at: string;
  last_seen_at: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  risk_score: number;
  risk_level: string;
  crypto_implementations?: CryptoImplementation[];
  highest_risk?: number;
}

export interface CryptoImplementation {
  id: string;
  tenant_id: string;
  asset_id: string;
  protocol: string;
  protocol_version?: string;
  cipher_suite?: string;
  key_exchange_algorithm?: string;
  signature_algorithm?: string;
  symmetric_encryption?: string;
  hash_algorithm?: string;
  key_size?: number;
  certificate_id?: string;
  discovery_method: string;
  confidence_score: number;
  source_sensor_id?: string;
  raw_data: Record<string, any>;
  risk_score: number;
  compliance_status: Record<string, any>;
  first_discovered_at: string;
  last_verified_at: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  risk_level: string;
  risk_factors?: string[];
}

export interface RiskSummary {
  total_assets: number;
  high_risk: number;
  medium_risk: number;
  low_risk: number;
  unknown_risk: number;
  total_crypto: number;
  critical_findings: number;
}

export interface AssetFilters {
  search?: string;
  asset_type?: string[];
  environment?: string[];
  risk_level?: string[];
  protocol?: string[];
  business_unit?: string[];
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_order?: string;
}

export interface AssetsResponse {
  assets: Asset[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}
