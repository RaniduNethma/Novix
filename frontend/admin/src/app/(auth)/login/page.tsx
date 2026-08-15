'use client';

import { useState } from 'react';
import { useAdminAuth } from '@/context/AdminAuthContext';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import toast from 'react-hot-toast';
import { Shield } from 'lucide-react';

export default function AdminLoginPage() {
  const { login } = useAdminAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await login(email, password);
      toast.success('Welcome back!');
    } catch (error: any) {
      toast.error(error.message || error.response?.data?.message || 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card>
      <div className="flex items-center gap-3 mb-6">
        <div className="p-2 bg-red-900/30 rounded-lg">
          <Shield className="h-5 w-5 text-red-400" />
        </div>
        <h2 className="text-xl font-bold text-white">Admin Sign In</h2>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Email"
          type="email"
          placeholder="admin@novix.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <Input
          label="Password"
          type="password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <Button type="submit" variant="primary" size="lg" isLoading={isLoading} className="w-full">
          Sign In
        </Button>
      </form>

      <p className="text-center text-xs text-gray-500 mt-4">Admin and Content Manager access only</p>
    </Card>
  );
}
