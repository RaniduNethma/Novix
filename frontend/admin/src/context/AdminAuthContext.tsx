'use client';

import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { User } from '@/types';
import { authService } from '@/services/auth.service';
import { useRouter } from 'next/navigation';

interface AdminAuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  isAdmin: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AdminAuthContext = createContext<AdminAuthContextType | undefined>(undefined);

export function AdminAuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const loadUser = async () => {
      if (authService.isAuthenticated()) {
        try {
          const userData = await authService.getMe();
          if (authService.isAdmin(userData)) {
            setUser(userData);
          } else {
            // Not an admin — reject
            authService.clearTokens();
            router.push('/login');
          }
        } catch {
          authService.clearTokens();
        }
      }
      setIsLoading(false);
    };
    loadUser();
  }, [router]);

  const login = useCallback(
    async (email: string, password: string) => {
      const response = await authService.login(email, password);

      if (!authService.isAdmin(response.user)) {
        throw new Error('Access denied. Admin privileges required.');
      }

      authService.saveTokens(response.accessToken, response.refreshToken);
      setUser(response.user);
      router.push('/');
    },
    [router],
  );

  const logout = useCallback(() => {
    authService.clearTokens();
    setUser(null);
    router.push('/login');
  }, [router]);

  return (
    <AdminAuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        isAdmin: !!user && authService.isAdmin(user),
        login,
        logout,
      }}
    >
      {children}
    </AdminAuthContext.Provider>
  );
}

export function useAdminAuth() {
  const context = useContext(AdminAuthContext);
  if (!context) {
    throw new Error('useAdminAuth must be used within AdminAuthProvider');
  }
  return context;
}
