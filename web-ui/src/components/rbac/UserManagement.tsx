/**
 * User Management Component
 * 
 * Provides comprehensive user management capabilities for tenant administrators.
 * This component handles user lifecycle management, role assignments, and user status.
 * 
 * Features:
 * - User listing with role assignments
 * - Role assignment and removal
 * - User activation/deactivation
 * - Role management modal
 * - User status tracking
 * 
 * Architecture:
 * - Mock data for demonstration (replace with API calls)
 * - Responsive table design with dark mode support
 * - Modal-based role assignment interface
 * - Real-time status updates
 */

import React, { useState, useEffect } from 'react';
import {
  UsersIcon,
  PencilIcon,
  ShieldCheckIcon,
  UserPlusIcon,
  XMarkIcon,
  CheckIcon
} from '@heroicons/react/24/outline';

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  roles: Role[];
  is_active: boolean;
  last_login_at?: string;
}

interface Role {
  id: string;
  name: string;
  display_name: string;
  description: string;
}

interface UserManagementProps {
  tenantId: string;
}

const UserManagement: React.FC<UserManagementProps> = ({ tenantId }) => {
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [showRoleAssignment, setShowRoleAssignment] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, [tenantId]);

  const loadData = async () => {
    try {
      setLoading(true);
      await Promise.all([
        loadUsers(),
        loadRoles()
      ]);
    } catch (error) {
      console.error('Failed to load user data:', error);
    } finally {
      setLoading(false);
    }
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
          { id: '1', name: 'tenant_admin', display_name: 'Tenant Administrator', description: 'Full tenant management' }
        ],
        is_active: true,
        last_login_at: '2024-01-15T10:30:00Z'
      },
      {
        id: '2',
        email: 'analyst@democorp.com',
        first_name: 'Security',
        last_name: 'Analyst',
        roles: [
          { id: '3', name: 'analyst', display_name: 'Security Analyst', description: 'Data analysis and reporting' }
        ],
        is_active: true,
        last_login_at: '2024-01-14T15:45:00Z'
      },
      {
        id: '3',
        email: 'viewer@democorp.com',
        first_name: 'Read',
        last_name: 'Only',
        roles: [
          { id: '4', name: 'viewer', display_name: 'Viewer', description: 'Read-only access' }
        ],
        is_active: true,
        last_login_at: '2024-01-13T09:20:00Z'
      }
    ];
    setUsers(mockUsers);
  };

  const loadRoles = async () => {
    // Mock data for demonstration
    const mockRoles: Role[] = [
      {
        id: '1',
        name: 'tenant_admin',
        display_name: 'Tenant Administrator',
        description: 'Full tenant management capabilities'
      },
      {
        id: '2',
        name: 'security_admin',
        display_name: 'Security Administrator',
        description: 'Security and compliance management'
      },
      {
        id: '3',
        name: 'analyst',
        display_name: 'Security Analyst',
        description: 'Data analysis and reporting'
      },
      {
        id: '4',
        name: 'viewer',
        display_name: 'Viewer',
        description: 'Read-only access'
      }
    ];
    setRoles(mockRoles);
  };

  const handleAssignRole = async (userId: string, roleId: string) => {
    try {
      // API call to assign role
      console.log(`Assigning role ${roleId} to user ${userId}`);
      // await assignUserRole(userId, roleId);
      loadData(); // Reload data
    } catch (error) {
      console.error('Failed to assign role:', error);
    }
  };

  const handleRemoveRole = async (userId: string, roleId: string) => {
    try {
      // API call to remove role
      console.log(`Removing role ${roleId} from user ${userId}`);
      // await removeUserRole(userId, roleId);
      loadData(); // Reload data
    } catch (error) {
      console.error('Failed to remove role:', error);
    }
  };

  const handleToggleUserStatus = async (userId: string, isActive: boolean) => {
    try {
      // API call to toggle user status
      console.log(`Setting user ${userId} active status to ${isActive}`);
      // await updateUserStatus(userId, isActive);
      loadData(); // Reload data
    } catch (error) {
      console.error('Failed to update user status:', error);
    }
  };

  const formatLastLogin = (lastLogin?: string) => {
    if (!lastLogin) return 'Never';
    const date = new Date(lastLogin);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };

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
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            User Management
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            Manage user accounts and role assignments
          </p>
        </div>
        <button
          onClick={() => console.log('Add User clicked')}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
        >
          <UserPlusIcon className="h-4 w-4 mr-2" />
          Add User
        </button>
      </div>

      {/* Users Table */}
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
                    Roles
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Last Login
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-600">
                {users.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="flex-shrink-0 h-10 w-10">
                          <div className="h-10 w-10 rounded-full bg-gray-300 dark:bg-gray-600 flex items-center justify-center">
                            <UsersIcon className="h-6 w-6 text-gray-600 dark:text-gray-300" />
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900 dark:text-white">
                            {user.first_name} {user.last_name}
                          </div>
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {user.email}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex flex-wrap gap-1">
                        {user.roles.map((role) => (
                          <span
                            key={role.id}
                            className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
                          >
                            {role.display_name}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          user.is_active
                            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                            : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                        }`}
                      >
                        {user.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {formatLastLogin(user.last_login_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => {
                            setSelectedUser(user);
                            setShowRoleAssignment(true);
                          }}
                          className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                          title="Manage Roles"
                        >
                          <ShieldCheckIcon className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => setSelectedUser(user)}
                          className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-300"
                          title="Edit User"
                        >
                          <PencilIcon className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleToggleUserStatus(user.id, !user.is_active)}
                          className={`${
                            user.is_active
                              ? 'text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300'
                              : 'text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300'
                          }`}
                          title={user.is_active ? 'Deactivate User' : 'Activate User'}
                        >
                          {user.is_active ? (
                            <XMarkIcon className="h-4 w-4" />
                          ) : (
                            <CheckIcon className="h-4 w-4" />
                          )}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Role Assignment Modal */}
      {showRoleAssignment && selectedUser && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white dark:bg-gray-800">
            <div className="mt-3">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">
                  Manage Roles for {selectedUser.first_name} {selectedUser.last_name}
                </h3>
                <button
                  onClick={() => setShowRoleAssignment(false)}
                  className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  <XMarkIcon className="h-6 w-6" />
                </button>
              </div>
              
              <div className="space-y-3">
                {roles.map((role) => {
                  const isAssigned = selectedUser.roles.some(r => r.id === role.id);
                  return (
                    <div
                      key={role.id}
                      className="flex items-center justify-between p-3 border rounded-lg dark:border-gray-600"
                    >
                      <div>
                        <h4 className="text-sm font-medium text-gray-900 dark:text-white">
                          {role.display_name}
                        </h4>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {role.description}
                        </p>
                      </div>
                      <button
                        onClick={() => {
                          if (isAssigned) {
                            handleRemoveRole(selectedUser.id, role.id);
                          } else {
                            handleAssignRole(selectedUser.id, role.id);
                          }
                        }}
                        className={`inline-flex items-center px-3 py-1 border text-sm font-medium rounded-md ${
                          isAssigned
                            ? 'border-red-300 text-red-700 bg-red-50 hover:bg-red-100 dark:border-red-600 dark:text-red-300 dark:bg-red-900/20 dark:hover:bg-red-900/30'
                            : 'border-blue-300 text-blue-700 bg-blue-50 hover:bg-blue-100 dark:border-blue-600 dark:text-blue-300 dark:bg-blue-900/20 dark:hover:bg-blue-900/30'
                        }`}
                      >
                        {isAssigned ? 'Remove' : 'Assign'}
                      </button>
                    </div>
                  );
                })}
              </div>
              
              <div className="flex justify-end mt-6">
                <button
                  onClick={() => setShowRoleAssignment(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 rounded-md hover:bg-gray-200 dark:hover:bg-gray-500"
                >
                  Close
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserManagement;
