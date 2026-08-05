import api from '@/lib/api';
import { AuthResponse, User } from '@/types';

export const authService = {
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/api/v1/auth/login', { email, password });
    return response.data;
  },

  async getMe(): Promise<User> {
    const response = await api.get<User>('/api/v1/users/me');
    return response.data;
  },

  saveTokens(accessToken: string, refreshToken: string) {
    localStorage.setItem('admin_accessToken', accessToken);
    localStorage.setItem('admin_refreshToken', refreshToken);
  },

  clearTokens() {
    localStorage.removeItem('admin_accessToken');
    localStorage.removeItem('admin_refreshToken');
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('admin_accessToken');
  },

  isAdmin(user: User): boolean {
    return user.roles?.includes('ROLE_ADMIN') || user.roles?.includes('ROLE_CONTENT_MANAGER');
  },
};
