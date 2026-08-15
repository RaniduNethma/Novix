'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAdminAuth } from '@/context/AdminAuthContext';
import { ADMIN_ROUTES } from '@/lib/constants';
import { clsx } from 'clsx';
import { LayoutDashboard, Video, Users, Tag, Bell, LogOut } from 'lucide-react';

const navItems = [
  {
    label: 'Dashboard',
    href: ADMIN_ROUTES.DASHBOARD,
    icon: LayoutDashboard,
  },
  {
    label: 'Videos',
    href: ADMIN_ROUTES.VIDEOS,
    icon: Video,
  },
  {
    label: 'Users',
    href: ADMIN_ROUTES.USERS,
    icon: Users,
  },
  {
    label: 'Categories',
    href: ADMIN_ROUTES.CATEGORIES,
    icon: Tag,
  },
  {
    label: 'Notifications',
    href: ADMIN_ROUTES.NOTIFICATIONS,
    icon: Bell,
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAdminAuth();

  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-gray-900 border-r border-gray-800 flex flex-col">
      {/* Logo */}
      <div className="p-6 border-b border-gray-800">
        <h1 className="text-xl font-bold text-red-500">NOVIX</h1>
        <p className="text-xs text-gray-500 mt-0.5">Admin Dashboard</p>
      </div>

      {/* Nav Items */}
      <nav className="flex-1 p-4 space-y-1">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={clsx(
                'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                isActive ? 'bg-red-600 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white',
              )}
            >
              <item.icon className="h-4 w-4" />
              {item.label}
            </Link>
          );
        })}
      </nav>

      {/* User + Logout */}
      <div className="p-4 border-t border-gray-800">
        <div className="flex items-center gap-3 mb-3">
          <div className="w-8 h-8 rounded-full bg-red-600 flex items-center justify-center text-sm font-bold">
            {user?.firstName?.[0]?.toUpperCase()}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-white truncate">
              {user?.firstName} {user?.lastName}
            </p>
            <p className="text-xs text-gray-500 truncate">{user?.email}</p>
          </div>
        </div>
        <button
          onClick={logout}
          className="flex items-center gap-2 w-full px-3 py-2 rounded-lg text-sm text-red-400 hover:bg-red-900/20 transition-colors"
        >
          <LogOut className="h-4 w-4" />
          Sign Out
        </button>
      </div>
    </aside>
  );
}
