import api from '@/lib/api';
import { PageResponse, Video } from '@/types';

export const videoService = {
  async getPublicVideos(page = 0, size = 20): Promise<PageResponse<Video>> {
    const response = await api.get<PageResponse<Video>>(`/api/v1/videos?page=${page}&size=${size}`);
    return response.data;
  },

  async getVideoById(videoId: string): Promise<Video> {
    const response = await api.get<Video>(`/api/v1/videos/${videoId}`);
    return response.data;
  },

  async getVideosByUploader(uploaderUserId: string, page = 0, size = 20): Promise<PageResponse<Video>> {
    const response = await api.get<PageResponse<Video>>(
      `/api/v1/videos/uploader/${uploaderUserId}` + `?page=${page}&size=${size}`,
    );
    return response.data;
  },
};
