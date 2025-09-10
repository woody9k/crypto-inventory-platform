import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { tenantsApi } from '../services/tenantsApi';
import {
  BuildingOfficeIcon,
  UsersIcon,
  ChartBarIcon,
  CurrencyDollarIcon,
} from '@heroicons/react/24/outline';

const DashboardPage: React.FC = () => {
  const { data: tenantsStats, isLoading } = useQuery({
    queryKey: ['tenants-stats'],
    queryFn: () => tenantsApi.getTenantsStats(),
  });

  const stats = [
    {
      name: 'Total Tenants',
      value: tenantsStats?.tenants_stats?.length || 0,
      icon: BuildingOfficeIcon,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      name: 'Active Tenants',
      value: tenantsStats?.tenants_stats?.filter(t => t.user_count > 0).length || 0,
      icon: ChartBarIcon,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
    {
      name: 'Total Users',
      value: tenantsStats?.tenants_stats?.reduce((sum, t) => sum + t.user_count, 0) || 0,
      icon: UsersIcon,
      color: 'text-purple-600',
      bgColor: 'bg-purple-100',
    },
    {
      name: 'Total Assets',
      value: tenantsStats?.tenants_stats?.reduce((sum, t) => sum + t.asset_count, 0) || 0,
      icon: CurrencyDollarIcon,
      color: 'text-yellow-600',
      bgColor: 'bg-yellow-100',
    },
  ];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-sm text-gray-500">
          Overview of your platform and tenant activity
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <div
              key={stat.name}
              className="relative bg-white pt-5 px-4 pb-12 sm:pt-6 sm:px-6 shadow rounded-lg overflow-hidden"
            >
              <dt>
                <div className={`absolute ${stat.bgColor} rounded-md p-3`}>
                  <Icon className={`h-6 w-6 ${stat.color}`} />
                </div>
                <p className="ml-16 text-sm font-medium text-gray-500 truncate">
                  {stat.name}
                </p>
              </dt>
              <dd className="ml-16 pb-6 flex items-baseline sm:pb-7">
                <p className="text-2xl font-semibold text-gray-900">{stat.value}</p>
              </dd>
            </div>
          );
        })}
      </div>

      {/* Recent Tenants */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
            Recent Tenants
          </h3>
          <div className="overflow-hidden">
            {tenantsStats?.tenants_stats?.length ? (
              <div className="space-y-4">
                {tenantsStats.tenants_stats.slice(0, 5).map((tenant) => (
                  <div
                    key={tenant.tenant_id}
                    className="flex items-center justify-between p-4 border border-gray-200 rounded-lg"
                  >
                    <div className="flex items-center">
                      <div className="flex-shrink-0">
                        <BuildingOfficeIcon className="h-8 w-8 text-gray-400" />
                      </div>
                      <div className="ml-4">
                        <h4 className="text-sm font-medium text-gray-900">
                          {tenant.tenant_name}
                        </h4>
                        <p className="text-sm text-gray-500">
                          {tenant.user_count} users • {tenant.asset_count} assets
                        </p>
                      </div>
                    </div>
                    <div className="text-sm text-gray-500">
                      {tenant.last_activity
                        ? new Date(tenant.last_activity).toLocaleDateString()
                        : 'No activity'}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-500 text-center py-8">No tenants found</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
