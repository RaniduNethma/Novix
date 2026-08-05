import api from '@/lib/api';

export const notificationService = {
  async sendNotification(userId: string, title: string, message: string, type = 'SYSTEM_ANNOUNCEMENT'): Promise<void> {
    await api.post('/api/v1/notifications/send', {
      userId,
      title,
      message,
      type,
    });
  },
};
