import api from '@/lib/api';
import { DashboardStats, PageResponse, Video } from '@/types';

export const statsService = {
  async getDashboardStats(): Promise<DashboardStats> {
    // Fetch data from multiple endpoints
    const [videosRes, categoriesRes] = await Promise.all([
      api.get<PageResponse<Video>>('/api/v1/videos?page=0&size=1'),
      api.get('/api/v1/categories?page=0&size=1'),
    ]);

    const totalVideos = videosRes.data.totalElements;
    const totalCategories = categoriesRes.data.totalElements;

    return {
      totalVideos,
      totalUsers: 0, // extend when user list endpoint added
      totalCategories,
      readyVideos: 0,
      processingVideos: 0,
      failedVideos: 0,
    };
  },
};
