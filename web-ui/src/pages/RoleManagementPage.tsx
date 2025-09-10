/**
 * Tenant Role Management Page
 * 
 * Main page for tenant-level RBAC (Role-Based Access Control) management.
 * This page provides a comprehensive interface for managing users, roles, and permissions
 * within the current tenant organization.
 * 
 * Features:
 * - Tab-based navigation for different management aspects
 * - Tenant role management with permission matrix
 * - Tenant user management with role assignments
 * - Permission overview by category
 * - Audit log viewing and tracking
 * 
 * Note: This is for tenant-level management only. Platform admin functions
 * are available in the separate SaaS Admin Console.
 * 
 * Architecture:
 * - PermissionProvider wrapper for RBAC context
 * - Tab-based interface with React state management
 * - Component composition for different management views
 * - Responsive design with dark mode support
 */

import React, { useState } from 'react';
import { PermissionProvider } from '../components/rbac/PermissionGate';
import RoleManagement from '../components/rbac/RoleManagement';
import UserManagement from '../components/rbac/UserManagement';
import {
  ShieldCheckIcon,
  UsersIcon,
  KeyIcon,
  ChartBarIcon
} from '@heroicons/react/24/outline';

// Permissions Overview Component
const PermissionsOverview: React.FC = () => {
  const permissions = [
    {
      category: 'Assets',
      permissions: [
        { name: 'assets.create', description: 'Create network assets' },
        { name: 'assets.read', description: 'View network assets' },
        { name: 'assets.update', description: 'Update network assets' },
        { name: 'assets.delete', description: 'Delete network assets' },
        { name: 'assets.manage', description: 'Full asset management' }
      ]
    },
    {
      category: 'Sensors',
      permissions: [
        { name: 'sensors.create', description: 'Create sensors' },
        { name: 'sensors.read', description: 'View sensors' },
        { name: 'sensors.update', description: 'Update sensors' },
        { name: 'sensors.delete', description: 'Delete sensors' },
        { name: 'sensors.manage', description: 'Full sensor management' }
      ]
    },
    {
      category: 'Reports',
      permissions: [
        { name: 'reports.create', description: 'Create reports' },
        { name: 'reports.read', description: 'View reports' },
        { name: 'reports.update', description: 'Update reports' },
        { name: 'reports.delete', description: 'Delete reports' },
        { name: 'reports.manage', description: 'Full report management' }
      ]
    },
    {
      category: 'Users',
      permissions: [
        { name: 'users.create', description: 'Create tenant users' },
        { name: 'users.read', description: 'View tenant users' },
        { name: 'users.update', description: 'Update tenant users' },
        { name: 'users.delete', description: 'Delete tenant users' },
        { name: 'users.manage', description: 'Full user management' }
      ]
    },
    {
      category: 'Settings',
      permissions: [
        { name: 'settings.read', description: 'View tenant settings' },
        { name: 'settings.update', description: 'Update tenant settings' },
        { name: 'settings.manage', description: 'Full settings management' }
      ]
    }
  ];

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
          Permission Overview
        </h2>
        <p className="text-gray-600 dark:text-gray-400">
          Available permissions in the system organized by category
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {permissions.map((category) => (
          <div key={category.category} className="bg-white dark:bg-gray-800 shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                {category.category}
              </h3>
              <div className="space-y-2">
                {category.permissions.map((permission) => (
                  <div
                    key={permission.name}
                    className="flex items-center justify-between p-3 border rounded-lg dark:border-gray-600"
                  >
                    <div>
                      <h4 className="text-sm font-medium text-gray-900 dark:text-white">
                        {permission.name}
                      </h4>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {permission.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

// Audit Logs Component
const AuditLogs: React.FC = () => {
  const auditLogs = [
    {
      id: '1',
      user: 'admin@democorp.com',
      action: 'permission_check',
      resource: 'assets',
      permission: 'assets.read',
      granted: true,
      timestamp: '2024-01-15T10:30:00Z',
      ip_address: '192.168.1.100'
    },
    {
      id: '2',
      user: 'analyst@democorp.com',
      action: 'role_assignment',
      resource: 'users',
      permission: 'users.read',
      granted: true,
      timestamp: '2024-01-15T09:45:00Z',
      ip_address: '192.168.1.101'
    },
    {
      id: '3',
      user: 'viewer@democorp.com',
      action: 'permission_check',
      resource: 'sensors',
      permission: 'sensors.manage',
      granted: false,
      timestamp: '2024-01-15T08:20:00Z',
      ip_address: '192.168.1.102'
    }
  ];

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
          Audit Logs
        </h2>
        <p className="text-gray-600 dark:text-gray-400">
          Track all permission checks and role changes
        </p>
      </div>

      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:p-6">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-600">
              <thead className="bg-gray-50 dark:bg-gray-700">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    User
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Action
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Resource
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Permission
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Result
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Timestamp
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-600">
                {auditLogs.map((log) => (
                  <tr key={log.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                      {log.user}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {log.action}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {log.resource}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {log.permission}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          log.granted
                            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                            : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                        }`}
                      >
                        {log.granted ? 'Granted' : 'Denied'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {formatTimestamp(log.timestamp)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

const RoleManagementPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'roles' | 'users' | 'permissions' | 'audit'>('roles');

  const tabs = [
    { id: 'roles', name: 'Roles', icon: ShieldCheckIcon },
    { id: 'users', name: 'Users', icon: UsersIcon },
    { id: 'permissions', name: 'Permissions', icon: KeyIcon },
    { id: 'audit', name: 'Audit Logs', icon: ChartBarIcon }
  ];

  const renderTabContent = () => {
    switch (activeTab) {
      case 'roles':
        return <RoleManagement />;
      case 'users':
        return <UserManagement tenantId="demo-tenant" />;
      case 'permissions':
        return <PermissionsOverview />;
      case 'audit':
        return <AuditLogs />;
      default:
        return <RoleManagement />;
    }
  };

  return (
    <PermissionProvider>
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          {/* Page Header */}
          <div className="mb-6">
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              Tenant Role Management
            </h1>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              Manage users, roles, and permissions for your organization. 
              <span className="text-blue-600 dark:text-blue-400">
                Platform admin functions are available in the separate SaaS Admin Console.
              </span>
            </p>
          </div>

          {/* Tab Navigation */}
          <div className="border-b border-gray-200 dark:border-gray-700 mb-6">
            <nav className="-mb-px flex space-x-8">
              {tabs.map((tab) => {
                const Icon = tab.icon;
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id as any)}
                    className={`${
                      activeTab === tab.id
                        ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
                    } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm flex items-center`}
                  >
                    <Icon className="h-4 w-4 mr-2" />
                    {tab.name}
                  </button>
                );
              })}
            </nav>
          </div>

          {/* Tab Content */}
          {renderTabContent()}
        </div>
      </div>
    </PermissionProvider>
  );
};

export default RoleManagementPage;
