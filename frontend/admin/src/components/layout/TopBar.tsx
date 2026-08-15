'use client';

import { usePathname } from 'next/navigation';

const pageTitles: Record<string, string> = {
  '/': 'Dashboard Overview',
  '/videos': 'Video Management',
  '/users': 'User Management',
  '/categories': 'Category Management',
  '/notifications': 'Notifications',
};

export function TopBar() {
  const pathname = usePathname();
  const title = pageTitles[pathname] || 'Admin';

  return (
    <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center px-6">
      <h2 className="text-lg font-semibold text-white">{title}</h2>
    </header>
  );
}
