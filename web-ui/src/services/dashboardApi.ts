// Dashboard API: aggregates lightweight overview data for the tenant dashboard by composing existing inventory endpoints.
import { inventoryApi } from './inventoryApi';

export interface DashboardOverview {
  total_assets: number;
  risk_summary: {
    total_assets: number;
    high_risk: number;
    medium_risk: number;
    low_risk: number;
    unknown_risk: number;
    total_crypto: number;
    critical_findings: number;
  } | null;
}

export const dashboardApi = {
  // Aggregate overview for the dashboard using existing endpoints.
  // Falls back gracefully if any sub-call fails.
  getOverview: async (): Promise<DashboardOverview> => {
    let totalAssets = 0;
    let riskSummary: DashboardOverview['risk_summary'] = null;

    try {
      const assetsResponse = await inventoryApi.getAssets({ page: 1, page_size: 1 });
      totalAssets = assetsResponse?.pagination?.total ?? 0;
    } catch (err) {
      totalAssets = 0;
    }

    try {
      const rs = await inventoryApi.getRiskSummary();
      riskSummary = rs ?? null;
    } catch (err) {
      riskSummary = null;
    }

    return {
      total_assets: totalAssets,
      risk_summary: riskSummary,
    };
  },
};
