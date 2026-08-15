import { Card } from '@/components/ui/Card';
import { clsx } from 'clsx';
import { LucideIcon } from 'lucide-react';

interface StatsCardProps {
  title: string;
  value: number | string;
  icon: LucideIcon;
  color?: 'red' | 'blue' | 'green' | 'yellow';
  subtitle?: string;
}

export function StatsCard({ title, value, icon: Icon, color = 'red', subtitle }: StatsCardProps) {
  return (
    <Card>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-gray-400">{title}</p>
          <p className="text-3xl font-bold text-white mt-1">
            {typeof value === 'number' ? value.toLocaleString() : value}
          </p>
          {subtitle && <p className="text-xs text-gray-500 mt-1">{subtitle}</p>}
        </div>
        <div
          className={clsx('p-3 rounded-xl', {
            'bg-red-900/30': color === 'red',
            'bg-blue-900/30': color === 'blue',
            'bg-green-900/30': color === 'green',
            'bg-yellow-900/30': color === 'yellow',
          })}
        >
          <Icon
            className={clsx('h-6 w-6', {
              'text-red-400': color === 'red',
              'text-blue-400': color === 'blue',
              'text-green-400': color === 'green',
              'text-yellow-400': color === 'yellow',
            })}
          />
        </div>
      </div>
    </Card>
  );
}
