'use client';

import { useState } from 'react';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { notificationService } from '@/services/notification.service';
import toast from 'react-hot-toast';
import { Send } from 'lucide-react';

export default function NotificationsPage() {
  const [form, setForm] = useState({
    userId: '',
    title: '',
    message: '',
    type: 'SYSTEM_ANNOUNCEMENT',
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await notificationService.sendNotification(form.userId, form.title, form.message, form.type);
      toast.success('Notification sent!');
      setForm({
        userId: '',
        title: '',
        message: '',
        type: 'SYSTEM_ANNOUNCEMENT',
      });
    } catch {
      toast.error('Failed to send notification');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-xl space-y-6">
      <Card>
        <h3 className="font-bold text-white mb-6">Send Notification</h3>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="User ID"
            placeholder="user123"
            value={form.userId}
            onChange={(e) =>
              setForm({
                ...form,
                userId: e.target.value,
              })
            }
            required
          />

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-gray-300">Type</label>
            <select
              value={form.type}
              onChange={(e) =>
                setForm({
                  ...form,
                  type: e.target.value,
                })
              }
              className="w-full px-3 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-red-500"
            >
              <option value="SYSTEM_ANNOUNCEMENT">System Announcement</option>
              <option value="VIDEO_PROCESSED">Video Processed</option>
              <option value="NEW_FOLLOWER">New Follower</option>
            </select>
          </div>

          <Input
            label="Title"
            placeholder="Notification title"
            value={form.title}
            onChange={(e) =>
              setForm({
                ...form,
                title: e.target.value,
              })
            }
            required
          />

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-gray-300">Message</label>
            <textarea
              value={form.message}
              onChange={(e) =>
                setForm({
                  ...form,
                  message: e.target.value,
                })
              }
              placeholder="Notification message..."
              rows={4}
              required
              className="w-full px-3 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 resize-none"
            />
          </div>

          <Button type="submit" variant="primary" size="lg" isLoading={isLoading} className="w-full">
            <Send className="h-4 w-4 mr-2" />
            Send Notification
          </Button>
        </form>
      </Card>
    </div>
  );
}
