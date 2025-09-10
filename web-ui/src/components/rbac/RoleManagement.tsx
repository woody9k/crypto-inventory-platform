/**
 * Role Management Component
 * 
 * Provides a comprehensive interface for managing user roles and permissions
 * in a multi-tenant SaaS environment. This component handles:
 * 
 * Features:
 * - Role creation and management
 * - Permission assignment and modification
 * - User role assignments
 * - Permission matrix visualization
 * - Real-time permission updates
 * 
 * Architecture:
 * - Mock data for demonstration (replace with API calls)
 * - Responsive design with dark mode support
 * - Tab-based interface for different management aspects
 * - Permission-based conditional rendering
 */

import React, { useState, useEffect } from 'react';
import {
  UsersIcon,
  ShieldCheckIcon,
  KeyIcon,
  PlusIcon,
  PencilIcon,
  CheckCircleIcon,
  XCircleIcon
} from '@heroicons/react/24/outline';

interface Role {
  id: string;
  name: string;
  display_name: string;
  description: string;
  is_system_role: boolean;
  permissions: Permission[];
}

interface Permission {
  id: string;
  name: string;
  resource: string;
  action: string;
  scope: string;
  description: string;
  granted: boolean;
}

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  roles: Role[];
}

const RoleManagement: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      // Load roles, users, and permissions
      // This would make API calls to the RBAC endpoints
      await Promise.all([
        loadRoles(),
        loadUsers(),
        loadPermissions()
      ]);
    } catch (error) {
      console.error('Failed to load RBAC data:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadRoles = async () => {
    // Mock data for demonstration
    const mockRoles: Role[] = [
      {
        id: '1',
        name: 'tenant_admin',
        display_name: 'Tenant Administrator',
        description: 'Full tenant management capabilities',
        is_system_role: true,
        permissions: []
      },
      {
        id: '2',
        name: 'security_admin',
        display_name: 'Security Administrator',
        description: 'Security and compliance management',
        is_system_role: true,
        permissions: []
      },
      {
        id: '3',
        name: 'analyst',
        display_name: 'Security Analyst',
        description: 'Data analysis and reporting',
        is_system_role: true,
        permissions: []
      },
      {
        id: '4',
        name: 'viewer',
        display_name: 'Viewer',
        description: 'Read-only access',
        is_system_role: true,
        permissions: []
      }
    ];
    setRoles(mockRoles);
  };

  const loadUsers = async () => {
    // Mock data for demonstration
    const mockUsers: User[] = [
      {
        id: '1',
        email: 'admin@democorp.com',
        first_name: 'Admin',
        last_name: 'User',
        roles: [
          {
            id: '1',
            name: 'tenant_admin',
            display_name: 'Tenant Administrator',
            description: 'Full tenant management capabilities',
            is_system_role: true,
            permissions: []
          }
        ]
      },
      {
        id: '2',
        email: 'analyst@democorp.com',
        first_name: 'Security',
        last_name: 'Analyst',
        roles: [
          {
            id: '3',
            name: 'analyst',
            display_name: 'Security Analyst',
            description: 'Data analysis and reporting',
            is_system_role: true,
            permissions: []
          }
        ]
      },
      {
        id: '3',
        email: 'viewer@democorp.com',
        first_name: 'Read',
        last_name: 'Only',
        roles: [
          {
            id: '4',
            name: 'viewer',
            display_name: 'Viewer',
            description: 'Read-only access',
            is_system_role: true,
            permissions: []
          }
        ]
      }
    ];
    setUsers(mockUsers);
  };

  const loadPermissions = async () => {
    // Mock data for demonstration
    const mockPermissions: Permission[] = [
      {
        id: '1',
        name: 'assets.create',
        resource: 'assets',
        action: 'create',
        scope: 'tenant',
        description: 'Create network assets',
        granted: false
      },
      {
        id: '2',
        name: 'assets.read',
        resource: 'assets',
        action: 'read',
        scope: 'tenant',
        description: 'View network assets',
        granted: true
      },
      {
        id: '3',
        name: 'sensors.manage',
        resource: 'sensors',
        action: 'manage',
        scope: 'tenant',
        description: 'Full sensor management',
        granted: false
      },
      {
        id: '4',
        name: 'users.manage',
        resource: 'users',
        action: 'manage',
        scope: 'tenant',
        description: 'Full user management',
        granted: true
      }
    ];
    setPermissions(mockPermissions);
  };

  const handleRolePermissionToggle = (permissionId: string, granted: boolean) => {
    if (selectedRole) {
      const updatedPermissions = permissions.map(p => 
        p.id === permissionId ? { ...p, granted } : p
      );
      setPermissions(updatedPermissions);
    }
  };

  // Placeholder functions for future API integration
  // const handleAssignRole = async (userId: string, roleId: string) => {
  //   try {
  //     // API call to assign role
  //     console.log(`Assigning role ${roleId} to user ${userId}`);
  //     // await assignUserRole(userId, roleId);
  //     loadData(); // Reload data
  //   } catch (error) {
  //     console.error('Failed to assign role:', error);
  //   }
  // };

  // const handleRemoveRole = async (userId: string, roleId: string) => {
  //   try {
  //     // API call to remove role
  //     console.log(`Removing role ${roleId} from user ${userId}`);
  //     // await removeUserRole(userId, roleId);
  //     loadData(); // Reload data
  //   } catch (error) {
  //     console.error('Failed to remove role:', error);
  //   }
  // };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Role-Based Access Control
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Manage user roles and permissions across your organization
          </p>
        </div>
        <div className="flex space-x-3">
          <button
            onClick={() => console.log('Create Role clicked')}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
          >
            <PlusIcon className="h-4 w-4 mr-2" />
            Create Role
          </button>
          <button
            onClick={() => console.log('Manage Users clicked')}
            className="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
          >
            <UsersIcon className="h-4 w-4 mr-2" />
            Manage Users
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <ShieldCheckIcon className="h-6 w-6 text-blue-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Total Roles
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {roles.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <UsersIcon className="h-6 w-6 text-green-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Total Users
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {users.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <KeyIcon className="h-6 w-6 text-purple-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Permissions
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {permissions.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <CheckCircleIcon className="h-6 w-6 text-green-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Active Users
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {users.filter(u => u.roles.length > 0).length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Roles and Users Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Roles Section */}
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
              Roles
            </h3>
            <div className="space-y-3">
              {roles.map((role) => (
                <div
                  key={role.id}
                  className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                    selectedRole?.id === role.id
                      ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                      : 'border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500'
                  }`}
                  onClick={() => setSelectedRole(role)}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm font-medium text-gray-900 dark:text-white">
                        {role.display_name}
                      </h4>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {role.description}
                      </p>
                    </div>
                    <div className="flex items-center space-x-2">
                      {role.is_system_role && (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                          System
                        </span>
                      )}
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setSelectedRole(role);
                        }}
                        className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                      >
                        <PencilIcon className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Users Section */}
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
              Users
            </h3>
            <div className="space-y-3">
              {users.map((user) => (
                <div
                  key={user.id}
                  className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                    selectedUser?.id === user.id
                      ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                      : 'border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500'
                  }`}
                  onClick={() => setSelectedUser(user)}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm font-medium text-gray-900 dark:text-white">
                        {user.first_name} {user.last_name}
                      </h4>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {user.email}
                      </p>
                      <div className="flex flex-wrap gap-1 mt-1">
                        {user.roles.map((role) => (
                          <span
                            key={role.id}
                            className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
                          >
                            {role.display_name}
                          </span>
                        ))}
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setSelectedUser(user);
                        }}
                        className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                      >
                        <PencilIcon className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Permission Matrix */}
      {selectedRole && (
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
              Permissions for {selectedRole.display_name}
            </h3>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-600">
                <thead className="bg-gray-50 dark:bg-gray-700">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Permission
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Resource
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Action
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Granted
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-600">
                  {permissions.map((permission) => (
                    <tr key={permission.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                        {permission.description}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {permission.resource}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {permission.action}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        <button
                          onClick={() => handleRolePermissionToggle(permission.id, !permission.granted)}
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            permission.granted
                              ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                              : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                          }`}
                        >
                          {permission.granted ? (
                            <CheckCircleIcon className="h-3 w-3 mr-1" />
                          ) : (
                            <XCircleIcon className="h-3 w-3 mr-1" />
                          )}
                          {permission.granted ? 'Granted' : 'Denied'}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default RoleManagement;
