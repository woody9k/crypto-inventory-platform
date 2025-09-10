import axios from 'axios';

const API_BASE_URL = 'http://localhost:8084/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('saas_admin_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('saas_admin_token');
      localStorage.removeItem('saas_admin_user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
