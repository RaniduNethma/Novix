import { clsx } from 'clsx';

type BadgeVariant = 'success' | 'warning' | 'danger' | 'info' | 'default';

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
}

export function Badge({ children, variant = 'default' }: BadgeProps) {
  return (
    <span
      className={clsx('inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium', {
        'bg-green-900/50 text-green-400': variant === 'success',
        'bg-yellow-900/50 text-yellow-400': variant === 'warning',
        'bg-red-900/50 text-red-400': variant === 'danger',
        'bg-blue-900/50 text-blue-400': variant === 'info',
        'bg-gray-800 text-gray-400': variant === 'default',
      })}
    >
      {children}
    </span>
  );
}
