import api from '@/lib/api';
import { Video, PageResponse, VideoStatus } from '@/types';

export const videoService = {
  async getVideos(page = 0, size = 20): Promise<PageResponse<Video>> {
    const response = await api.get<PageResponse<Video>>(`/api/v1/videos?page=${page}&size=${size}`);
    return response.data;
  },

  async getVideoById(videoId: string): Promise<Video> {
    const response = await api.get<Video>(`/api/v1/videos/${videoId}`);
    return response.data;
  },

  async updateVideoStatus(videoId: string, status: VideoStatus): Promise<void> {
    await api.patch(`/api/v1/videos/${videoId}/status?status=${status}`);
  },

  async deleteVideo(videoId: string): Promise<void> {
    await api.delete(`/api/v1/videos/${videoId}`);
  },
};
