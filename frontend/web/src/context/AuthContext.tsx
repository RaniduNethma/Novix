'use client';

import { authService } from '@/services/auth.service';
import { LoginRequest, RegisterRequest, User } from '@/types';
import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  // Load user on mount
  useEffect(() => {
    const loadUser = async () => {
      if (authService.isAuthenticated()) {
        try {
          const userData = await authService.getMe();
          setUser(userData);
        } catch {
          authService.clearTokens();
        }
      }
      setIsLoading(false);
    };
    loadUser();
  }, []);

  const login = useCallback(
    async (data: LoginRequest) => {
      const response = await authService.login(data);
      authService.saveTokens(response.accessToken, response.refreshToken);
      setUser(response.user);
      router.push('/');
    },
    [router],
  );

  const register = useCallback(
    async (data: RegisterRequest) => {
      const response = await authService.register(data);
      authService.saveTokens(response.accessToken, response.refreshToken);
      setUser(response.user);
      router.push('/');
    },
    [router],
  );

  const logout = useCallback(async () => {
    await authService.logout();
    setUser(null);
    router.push('/login');
  }, [router]);

  return (
    <AuthContext.Provider value={{ user, isLoading, isAuthenticated: !!user, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
