/**
 * Permission Gate Component
 * 
 * Provides a comprehensive permission-based access control system for React components.
 * This system allows for granular control over what users can see and interact with
 * based on their assigned roles and permissions.
 * 
 * Features:
 * - Context-based permission management
 * - Permission gates for conditional rendering
 * - Higher-order components for permission-based wrapping
 * - Hooks for permission checking
 * - Support for multiple permission types (any/all)
 * 
 * Usage:
 * - Wrap your app with <PermissionProvider>
 * - Use <PermissionGate> to conditionally render components
 * - Use usePermissions() hook for programmatic permission checks
 */

import React, { createContext, useContext, useState, useEffect } from 'react';

interface PermissionContextType {
  permissions: string[];
  hasPermission: (permission: string) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasAllPermissions: (permissions: string[]) => boolean;
  loading: boolean;
}

const PermissionContext = createContext<PermissionContextType | undefined>(undefined);

export const usePermissions = () => {
  const context = useContext(PermissionContext);
  if (context === undefined) {
    throw new Error('usePermissions must be used within a PermissionProvider');
  }
  return context;
};

interface PermissionProviderProps {
  children: React.ReactNode;
}

export const PermissionProvider: React.FC<PermissionProviderProps> = ({ children }) => {
  const [permissions, setPermissions] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadUserPermissions();
  }, []);

  const loadUserPermissions = async () => {
    try {
      // This would make an API call to get the current user's permissions
      // For now, we'll use mock data
      const mockPermissions = [
        'assets.read',
        'sensors.read',
        'reports.read',
        'users.read',
        'settings.read'
      ];
      setPermissions(mockPermissions);
    } catch (error) {
      console.error('Failed to load user permissions:', error);
    } finally {
      setLoading(false);
    }
  };

  const hasPermission = (permission: string): boolean => {
    return permissions.includes(permission);
  };

  const hasAnyPermission = (requiredPermissions: string[]): boolean => {
    return requiredPermissions.some(permission => permissions.includes(permission));
  };

  const hasAllPermissions = (requiredPermissions: string[]): boolean => {
    return requiredPermissions.every(permission => permissions.includes(permission));
  };

  const value: PermissionContextType = {
    permissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    loading
  };

  return (
    <PermissionContext.Provider value={value}>
      {children}
    </PermissionContext.Provider>
  );
};

interface PermissionGateProps {
  permission?: string;
  permissions?: string[];
  requireAll?: boolean;
  children: React.ReactNode;
  fallback?: React.ReactNode;
  loadingFallback?: React.ReactNode;
}

export const PermissionGate: React.FC<PermissionGateProps> = ({
  permission,
  permissions = [],
  requireAll = false,
  children,
  fallback = null,
  loadingFallback = null
}) => {
  const { hasAnyPermission, hasAllPermissions, loading } = usePermissions();

  if (loading) {
    return <>{loadingFallback}</>;
  }

  let hasAccess = false;

  if (permission) {
    hasAccess = hasAnyPermission([permission]);
  } else if (permissions.length > 0) {
    hasAccess = requireAll ? hasAllPermissions(permissions) : hasAnyPermission(permissions);
  } else {
    // If no permission specified, allow access
    hasAccess = true;
  }

  return hasAccess ? <>{children}</> : <>{fallback}</>;
};

// Higher-order component for permission-based rendering
export const withPermission = <P extends object>(
  Component: React.ComponentType<P>,
  permission: string,
  fallback?: React.ComponentType<P>
) => {
  return (props: P) => (
    <PermissionGate permission={permission} fallback={fallback ? React.createElement(fallback, props) : null}>
      <Component {...props} />
    </PermissionGate>
  );
};

// Hook for conditional rendering based on permissions
export const usePermissionCheck = (permission: string) => {
  const { hasAnyPermission, loading } = usePermissions();
  return {
    hasPermission: hasAnyPermission([permission]),
    loading
  };
};

// Hook for multiple permission checks
export const usePermissionChecks = (permissions: string[], requireAll = false) => {
  const { hasAnyPermission, hasAllPermissions, loading } = usePermissions();
  
  return {
    hasPermission: requireAll ? hasAllPermissions(permissions) : hasAnyPermission(permissions),
    hasAllPermissions: hasAllPermissions(permissions),
    hasAnyPermission: hasAnyPermission(permissions),
    loading
  };
};
