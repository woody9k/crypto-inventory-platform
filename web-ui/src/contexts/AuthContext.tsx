import React, { createContext, useContext, useEffect, useState } from 'react';
import { AuthContextType, User, Tenant, LoginRequest, RegisterRequest } from '../types';
import { authApi, tokenManager } from '../services/api';
import toast from 'react-hot-toast';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user && !!tokenManager.getToken();

  // Initialize auth state from existing token
  useEffect(() => {
    const initializeAuth = async () => {
      const token = tokenManager.getToken();
      if (token) {
        try {
          // Try to get current user info
          const userData = await authApi.getCurrentUser();
          setUser(userData.user);
          if (userData.tenant) {
            setTenant(userData.tenant);
          }
        } catch (error) {
          console.error('Failed to initialize auth:', error);
          // Token might be expired, clear it
          tokenManager.clearTokens();
        }
      }
      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    try {
      setIsLoading(true);
      const response = await authApi.login(credentials);
      
      // Store tokens
      tokenManager.setToken(response.access_token);
      tokenManager.setRefreshToken(response.refresh_token);
      
      // Set user data
      setUser(response.user);
      
      // Note: Tenant data might come from user object or separate API call
      // For now, we'll extract it if it's embedded in the user response
      if ('tenant' in response && response.tenant) {
        setTenant(response.tenant as any);
      }
      
      toast.success('Successfully logged in!');
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || error.response?.data?.details || 'Login failed';
      toast.error(errorMessage);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (data: RegisterRequest): Promise<void> => {
    try {
      setIsLoading(true);
      const response = await authApi.register(data);
      
      if (response.error) {
        throw new Error(response.error);
      }
      
      toast.success('Registration successful! Please log in.');
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || error.response?.data?.details || 'Registration failed';
      toast.error(errorMessage);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async (): Promise<void> => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // Always clear local state
      setUser(null);
      setTenant(null);
      tokenManager.clearTokens();
      toast.success('Logged out successfully');
    }
  };

  const refreshAuth = async (): Promise<void> => {
    const refreshToken = tokenManager.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await authApi.refreshToken(refreshToken);
      tokenManager.setToken(response.access_token);
      tokenManager.setRefreshToken(response.refresh_token);
      setUser(response.user);
    } catch (error) {
      // Refresh failed, clear tokens and logout
      logout();
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    tenant,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    refreshAuth,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
