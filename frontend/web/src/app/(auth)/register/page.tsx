'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { ROUTES } from '@/lib/constants';
import toast from 'react-hot-toast';

export default function RegisterPage() {
  const { register } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [form, setForm] = useState({
    firstName: '',
    lastName: '',
    username: '',
    email: '',
    password: '',
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validate = () => {
    const newErrors: Record<string, string> = {};
    if (!form.firstName) newErrors.firstName = 'First name is required';
    if (!form.lastName) newErrors.lastName = 'Last name is required';
    if (!form.username || form.username.length < 3) newErrors.username = 'Username must be at least 3 characters';
    if (!form.email) newErrors.email = 'Email is required';
    if (!form.password || form.password.length < 8) newErrors.password = 'Password must be at least 8 characters';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setIsLoading(true);
    try {
      await register(form);
      toast.success('Account created successfully!');
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Registration failed');
    } finally {
      setIsLoading(false);
    }
  };

  const updateForm = (field: string, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  return (
    <Card>
      <h2 className="text-xl font-bold text-white mb-6">Create your account</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <Input
            label="First Name"
            placeholder="John"
            value={form.firstName}
            onChange={(e) => updateForm('firstName', e.target.value)}
            error={errors.firstName}
          />
          <Input
            label="Last Name"
            placeholder="Doe"
            value={form.lastName}
            onChange={(e) => updateForm('lastName', e.target.value)}
            error={errors.lastName}
          />
        </div>

        <Input
          label="Username"
          placeholder="johndoe"
          value={form.username}
          onChange={(e) => updateForm('username', e.target.value)}
          error={errors.username}
        />

        <Input
          label="Email"
          type="email"
          placeholder="you@example.com"
          value={form.email}
          onChange={(e) => updateForm('email', e.target.value)}
          error={errors.email}
        />

        <Input
          label="Password"
          type="password"
          placeholder="••••••••"
          value={form.password}
          onChange={(e) => updateForm('password', e.target.value)}
          error={errors.password}
        />

        <Button type="submit" variant="primary" size="lg" isLoading={isLoading} className="w-full">
          Create Account
        </Button>
      </form>

      <p className="text-center text-sm text-gray-400 mt-6">
        Already have an account?{' '}
        <Link href={ROUTES.LOGIN} className="text-red-400 hover:text-red-300">
          Sign in
        </Link>
      </p>
    </Card>
  );
}
