import { clsx } from 'clsx';

export function Card({ children, className }: { children: React.ReactNode; className?: string }) {
  return <div className={clsx('bg-gray-900 rounded-xl border border-gray-800 p-6', className)}>{children}</div>;
}
