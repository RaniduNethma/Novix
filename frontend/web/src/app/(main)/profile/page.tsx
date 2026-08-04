'use client';

import { useAuth } from '@/context/AuthContext';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { ROUTES } from '@/lib/constants';
import { User, Mail, Shield, Calendar } from 'lucide-react';

export default function ProfilePage() {
  const { user, isAuthenticated, isLoading, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push(ROUTES.LOGIN);
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading || !user) {
    return null;
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-white">My Profile</h1>

      <Card>
        <div className="flex items-center gap-4 mb-6">
          <div className="w-16 h-16 rounded-full bg-red-600 flex items-center justify-center text-2xl font-bold">
            {user.firstName?.[0]?.toUpperCase()}
          </div>
          <div>
            <h2 className="text-xl font-bold text-white">
              {user.firstName} {user.lastName}
            </h2>
            <p className="text-gray-400">@{user.username}</p>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-3 text-gray-300">
            <Mail className="h-4 w-4 text-gray-500" />
            <span>{user.email}</span>
          </div>

          <div className="flex items-center gap-3 text-gray-300">
            <Shield className="h-4 w-4 text-gray-500" />
            <div className="flex gap-2">
              {user.roles?.map((role) => (
                <span key={role} className="bg-gray-800 text-xs px-2 py-0.5 rounded-full">
                  {role.replace('ROLE_', '')}
                </span>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-3 text-gray-300">
            <Calendar className="h-4 w-4 text-gray-500" />
            <span className="text-sm">
              Joined{' '}
              {new Date(user.createdAt).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
              })}
            </span>
          </div>
        </div>
      </Card>

      <Card>
        <h3 className="font-bold text-white mb-4">Account Actions</h3>
        <Button variant="danger" onClick={logout}>
          Sign Out
        </Button>
      </Card>
    </div>
  );
}
