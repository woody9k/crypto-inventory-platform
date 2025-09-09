/**
 * Main Application Component
 * 
 * This is the root component of the crypto inventory management platform.
 * It sets up the application routing, state management, and global providers.
 * 
 * Architecture:
 * - React Router v6 for client-side routing with protected routes
 * - React Query for server state management and caching
 * - Context providers for authentication and theme management
 * - Toast notifications for user feedback
 * - Responsive layout with header navigation
 * 
 * Key Features:
 * - Multi-tenant authentication with JWT tokens
 * - Protected routing based on authentication status
 * - Dark/light theme support with persistent preferences
 * - Real-time data fetching with React Query
 * - Professional UI with TailwindCSS styling
 */
import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { Header } from './components/layout/Header';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { DashboardPage } from './pages/DashboardPage';
import { AssetsPage } from './pages/AssetsPage';
import SensorManagementPage from './pages/SensorManagementPage';
import SensorRegistrationPage from './pages/SensorRegistrationPage';
import ReportsPage from './pages/ReportsPage';
import RoleManagementPage from './pages/RoleManagementPage';

// Create a query client instance with optimized defaults for crypto inventory data
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

// Protected Route Component
interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

// Layout component that includes header
const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Header />
      <main>{children}</main>
    </div>
  );
};

// App Router Component
const AppRouter: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route
          path="/roles"
          element={
            <ProtectedRoute>
              <Layout>
                <RoleManagementPage />
              </Layout>
            </ProtectedRoute>
          }
        />
        {/* Public routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        
        {/* Default route - redirect to dashboard */}
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        
        {/* Protected routes */}
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Layout>
                <DashboardPage />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/assets"
          element={
            <ProtectedRoute>
              <Layout>
                <AssetsPage />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/sensors"
          element={
            <ProtectedRoute>
              <Layout>
                <SensorManagementPage />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/sensors/register"
          element={
            <ProtectedRoute>
              <Layout>
                <SensorRegistrationPage />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/reports"
          element={
            <ProtectedRoute>
              <Layout>
                <ReportsPage />
              </Layout>
            </ProtectedRoute>
          }
        />
        {/* Default redirects */}
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Router>
  );
};

// Main App Component
const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AuthProvider>
          <AppRouter />
          <Toaster
            position="top-right"
            toastOptions={{
              duration: 4000,
              style: {
                background: 'var(--toast-bg)',
                color: 'var(--toast-color)',
              },
              className: 'dark:bg-gray-800 dark:text-white',
            }}
          />
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
};

export default App;
