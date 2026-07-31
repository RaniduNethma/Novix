import api from '@/lib/api';
import { AuthResponse, LoginRequest, RegisterRequest, User } from '@/types';

export const authService = {
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/api/v1/auth/register', data);
    return response.data;
  },

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/api/v1/auth/register', data);
    return response.data;
  },

  async getMe(): Promise<User> {
    const response = await api.get<User>('/api/v1/users/me');
    return response.data;
  },

  async logout(): Promise<void> {
    const token = localStorage.getItem('accessToken');
    await api.post('api/v1/auth/logout', {}, { headers: { Authorization: `Bearer ${token}` } });
    localStorage.clear();
  },

  saveTokens(accessToken: string, refreshToken: string) {
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('refreshToken', refreshToken);
  },

  clearTokens() {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('accessToken');
  },
};
