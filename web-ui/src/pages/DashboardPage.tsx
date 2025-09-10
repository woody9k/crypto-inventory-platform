import React, { useMemo, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { 
  ChartBarIcon,
  ServerIcon,
  ShieldCheckIcon,
  DocumentTextIcon,
  CpuChipIcon,
  UserGroupIcon
} from '@heroicons/react/24/outline';
import { useQuery } from '@tanstack/react-query';
import { dashboardApi } from '../services/dashboardApi';
import { Link, useNavigate } from 'react-router-dom';
import { auditApi } from '../services/auditApi';
import { certificatesApi } from '../services/certificatesApi';
import { inventoryApi } from '../services/inventoryApi';

const rangeOptions = [
  { label: '24h', days: 1 },
  { label: '7d', days: 7 },
  { label: '30d', days: 30 },
];

export const DashboardPage: React.FC = () => {
  const { user, tenant } = useAuth();
  const navigate = useNavigate();
  const [range, setRange] = useState(rangeOptions[1]);

  const { fromISO, toISO } = useMemo(() => {
    const to = new Date();
    const from = new Date();
    from.setDate(to.getDate() - range.days);
    return { fromISO: from.toISOString(), toISO: to.toISOString() };
  }, [range]);

  const { data: overview, isLoading, isError } = useQuery({
    queryKey: ['dashboard-overview', tenant?.id],
    queryFn: () => dashboardApi.getOverview(),
    enabled: true,
  });

  const { data: audit, isLoading: isAuditLoading, isError: isAuditError } = useQuery({
    queryKey: ['audit-logs', tenant?.id, fromISO, toISO],
    queryFn: () => auditApi.getLogs({ tenant_id: tenant?.id, limit: 6, page: 1, from: fromISO, to: toISO }),
    enabled: !!tenant?.id,
  });

  const { data: expiringCerts, isLoading: isCertsLoading, isError: isCertsError } = useQuery({
    queryKey: ['expiring-certs', tenant?.id],
    queryFn: () => certificatesApi.getExpiring(5),
    enabled: !!tenant?.id,
  });

  const { data: riskSummary, isLoading: isRiskLoading, isError: isRiskError } = useQuery({
    queryKey: ['risk-summary', tenant?.id],
    queryFn: () => inventoryApi.getRiskSummary(),
    enabled: !!tenant?.id,
  });

  const assetsCount = overview?.total_assets ?? 0;
  const cryptoTotal = overview?.risk_summary?.total_crypto ?? 0;

  const stats = [
    { name: 'Network Assets', value: assetsCount.toLocaleString(), icon: ServerIcon, change: '', changeType: 'neutral' as const },
    { name: 'Crypto Implementations', value: cryptoTotal.toLocaleString(), icon: CpuChipIcon, change: '', changeType: 'neutral' as const },
    { name: 'Compliance Score', value: '—', icon: ShieldCheckIcon, change: '', changeType: 'neutral' as const },
    { name: 'Active Sensors', value: '—', icon: ChartBarIcon, change: '', changeType: 'neutral' as const },
  ];

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header with Range Selector */}
        <div className="mb-8 flex items-start sm:items-center justify-between flex-col sm:flex-row gap-4">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
              Welcome back, {user?.first_name}!
            </h1>
            <p className="mt-2 text-gray-600 dark:text-gray-400">
              Here's what's happening with your crypto inventory at {tenant?.name || 'your organization'}.
            </p>
          </div>
          <div className="flex items-center space-x-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-1">
            {rangeOptions.map((opt) => (
              <button
                key={opt.label}
                className={`px-3 py-1 text-sm rounded-md ${opt.label === range.label ? 'bg-primary-600 text-white' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'}`}
                onClick={() => setRange(opt)}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {isLoading && Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
              <div className="p-5 animate-pulse">
                <div className="h-6 w-6 bg-gray-200 dark:bg-gray-700 rounded" />
                <div className="mt-4 h-4 w-24 bg-gray-200 dark:bg-gray-700 rounded" />
                <div className="mt-2 h-6 w-20 bg-gray-200 dark:bg-gray-700 rounded" />
              </div>
            </div>
          ))}

          {!isLoading && !isError && stats.map((stat) => (
            <div
              key={stat.name}
              className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg"
            >
              <div className="p-5">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <stat.icon className="h-6 w-6 text-gray-400" aria-hidden="true" />
                  </div>
                  <div className="ml-5 w-0 flex-1">
                    <dl>
                      <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                        {stat.name}
                      </dt>
                      <dd className="flex items-baseline">
                        <div className="text-2xl font-semibold text-gray-900 dark:text-white">
                          {stat.value}
                        </div>
                      </dd>
                    </dl>
                  </div>
                </div>
              </div>
            </div>
          ))}

          {!isLoading && isError && (
            <div className="md:col-span-2 lg:col-span-4 bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg p-5">
              <p className="text-sm text-red-600">Failed to load dashboard stats. Please try again.</p>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Recent Activity */}
          <div className="bg-white dark:bg-gray-800 shadow rounded-lg lg:col-span-2">
            <div className="px-4 py-5 sm:p-6">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
                  Recent Activity
                </h3>
                <Link to="/reports" className="text-sm text-primary-600 hover:text-primary-500">View reports</Link>
              </div>

              {isAuditLoading && (
                <div className="space-y-3 animate-pulse">
                  {Array.from({ length: 4 }).map((_, i) => (
                    <div key={i} className="h-6 bg-gray-100 dark:bg-gray-700 rounded" />
                  ))}
                </div>
              )}

              {!isAuditLoading && isAuditError && (
                <p className="text-sm text-red-600">Failed to load recent activity.</p>
              )}

              {!isAuditLoading && !isAuditError && (
                <div className="space-y-3">
                  {(audit?.logs?.length ?? 0) === 0 && (
                    <p className="text-sm text-gray-500 dark:text-gray-400">No recent activity.</p>
                  )}
                  {audit?.logs?.map((log) => (
                    <div key={log.id} className="flex items-center space-x-3">
                      <div className="flex-shrink-0">
                        <div className="w-2 h-2 bg-primary-500 rounded-full"></div>
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {log.action} on {log.resource_type}
                        </p>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {new Date(log.created_at).toLocaleString()}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Expiring Certificates */}
          <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">Expiring Certificates</h3>

              {isCertsLoading && (
                <div className="space-y-3 animate-pulse">
                  {Array.from({ length: 5 }).map((_, i) => (
                    <div key={i} className="h-6 bg-gray-100 dark:bg-gray-700 rounded" />
                  ))}
                </div>
              )}

              {!isCertsLoading && isCertsError && (
                <p className="text-sm text-red-600">Failed to load expiring certificates.</p>
              )}

              {!isCertsLoading && !isCertsError && (
                <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                  {(expiringCerts?.length ?? 0) === 0 && (
                    <p className="text-sm text-gray-500 dark:text-gray-400">No certificates nearing expiry.</p>
                  )}
                  {expiringCerts?.map((cert) => (
                    <li key={cert.id} className="py-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium text-gray-900 dark:text-white">{cert.common_name || cert.id}</p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">{cert.issuer || 'Unknown issuer'}</p>
                        </div>
                        <div className="text-sm text-gray-700 dark:text-gray-300">
                          {cert.days_until_expiry} days
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              )}
              <div className="mt-4">
                <button onClick={() => navigate('/assets')} className="text-sm text-primary-600 hover:text-primary-500">View all certificates</button>
              </div>
            </div>
          </div>
        </div>

        {/* Top Risks */}
        <div className="mt-6 bg-white dark:bg-gray-800 shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">Top Risks</h3>
            {isRiskLoading && (
              <div className="space-y-3 animate-pulse">
                {Array.from({ length: 3 }).map((_, i) => (
                  <div key={i} className="h-6 bg-gray-100 dark:bg-gray-700 rounded" />
                ))}
              </div>
            )}
            {!isRiskLoading && isRiskError && (
              <p className="text-sm text-red-600">Failed to load risk summary.</p>
            )}
            {!isRiskLoading && !isRiskError && riskSummary && (
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
                  <p className="text-xs text-gray-500 dark:text-gray-400">High Risk</p>
                  <p className="text-xl font-semibold text-red-600">{riskSummary.high_risk}</p>
                </div>
                <div className="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
                  <p className="text-xs text-gray-500 dark:text-gray-400">Medium Risk</p>
                  <p className="text-xl font-semibold text-yellow-600">{riskSummary.medium_risk}</p>
                </div>
                <div className="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
                  <p className="text-xs text-gray-500 dark:text-gray-400">Low Risk</p>
                  <p className="text-xl font-semibold text-green-600">{riskSummary.low_risk}</p>
                </div>
                <div className="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
                  <p className="text-xs text-gray-500 dark:text-gray-400">Unknown</p>
                  <p className="text-xl font-semibold text-gray-600 dark:text-gray-300">{riskSummary.unknown_risk}</p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Quick Actions */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
          <button className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 text-left w-full" onClick={() => navigate('/assets')}>
            <div className="flex items-center">
              <ServerIcon className="h-8 w-8 text-primary-600" />
              <div className="ml-4">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">
                  Asset Discovery
                </h3>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  Scan your network for crypto implementations
                </p>
              </div>
            </div>
          </button>

          <button className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 text-left w-full" onClick={() => navigate('/reports')}>
            <div className="flex items-center">
              <DocumentTextIcon className="h-8 w-8 text-primary-600" />
              <div className="ml-4">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">
                  Generate Report
                </h3>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  Create compliance and inventory reports
                </p>
              </div>
            </div>
          </button>

          <button className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 text-left w-full" onClick={() => navigate('/roles')}>
            <div className="flex items-center">
              <UserGroupIcon className="h-8 w-8 text-primary-600" />
              <div className="ml-4">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">
                  Manage Team
                </h3>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  Add team members and manage permissions
                </p>
              </div>
            </div>
          </button>
        </div>
      </div>
    </div>
  );
};
