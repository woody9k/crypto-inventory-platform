import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import { 
  ChartBarIcon,
  ServerIcon,
  ShieldCheckIcon,
  DocumentTextIcon,
  CpuChipIcon,
  UserGroupIcon
} from '@heroicons/react/24/outline';

export const DashboardPage: React.FC = () => {
  const { user, tenant } = useAuth();

  const stats = [
    { name: 'Network Assets', value: '1,247', icon: ServerIcon, change: '+12%', changeType: 'positive' },
    { name: 'Crypto Implementations', value: '892', icon: CpuChipIcon, change: '+7%', changeType: 'positive' },
    { name: 'Compliance Score', value: '78%', icon: ShieldCheckIcon, change: '+5%', changeType: 'positive' },
    { name: 'Active Sensors', value: '12', icon: ChartBarIcon, change: '+2', changeType: 'positive' },
  ];

  const recentActivity = [
    { action: 'New TLS 1.3 implementation detected', asset: 'web-prod-01.democorp.com', time: '2 hours ago' },
    { action: 'Certificate expiring soon', asset: 'db-primary.democorp.internal', time: '4 hours ago' },
    { action: 'Sensor heartbeat received', asset: 'datacenter-sensor-01', time: '5 minutes ago' },
    { action: 'Compliance check completed', asset: 'PCI DSS Assessment', time: '1 day ago' },
    { action: 'Report generated', asset: 'Crypto Summary Report', time: '3 hours ago' },
    { action: 'Risk assessment updated', asset: 'Network Security Scan', time: '6 hours ago' },
  ];

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Welcome Section */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            Welcome back, {user?.first_name}!
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Here's what's happening with your crypto inventory at {tenant?.name || 'your organization'}.
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {stats.map((stat) => (
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
                        <div
                          className={`ml-2 flex items-baseline text-sm font-semibold ${
                            stat.changeType === 'positive'
                              ? 'text-green-600'
                              : stat.changeType === 'negative'
                              ? 'text-red-600'
                              : 'text-gray-500'
                          }`}
                        >
                          {stat.change}
                        </div>
                      </dd>
                    </dl>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Recent Activity */}
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">
              Recent Activity
            </h3>
            <div className="space-y-3">
              {recentActivity.map((activity, index) => (
                <div key={index} className="flex items-center space-x-3">
                  <div className="flex-shrink-0">
                    <div className="w-2 h-2 bg-primary-500 rounded-full"></div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 dark:text-white">
                      {activity.action}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {activity.asset}
                    </p>
                  </div>
                  <div className="flex-shrink-0 text-sm text-gray-500 dark:text-gray-400">
                    {activity.time}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
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
          </div>

          <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
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
          </div>

          <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
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
          </div>
        </div>
      </div>
    </div>
  );
};
