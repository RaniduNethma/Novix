'use client';

import { Card } from '@/components/ui/Card';
import { Users } from 'lucide-react';

export default function UsersPage() {
  return (
    <Card>
      <div className="text-center py-16">
        <div className="flex justify-center mb-4">
          <div className="p-4 bg-gray-800 rounded-full">
            <Users className="h-8 w-8 text-gray-500" />
          </div>
        </div>
        <h3 className="text-lg font-bold text-white mb-2">User Management</h3>
        <p className="text-gray-400 text-sm max-w-sm mx-auto">
          User list endpoint coming soon. Add a paginated user list endpoint to user-service to enable this feature.
        </p>
      </div>
    </Card>
  );
}
