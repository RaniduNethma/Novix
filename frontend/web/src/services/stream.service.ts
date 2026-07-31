import api from '@/lib/api';
import { API_ENDPOINTS } from '@/lib/constants';
import { StreamResponse } from '@/types';

export const streamService = {
  async startStream(videoId: string): Promise<StreamResponse> {
    const response = await api.post<StreamResponse>(API_ENDPOINTS.STREAM_START, { videoId });
    return response.data;
  },

  async endStream(sessionId: string): Promise<void> {
    await api.delete(API_ENDPOINTS.STREAM_END(sessionId));
  },

  getManifestUrl(videoId: string): string {
    return `${process.env.NEXT_PUBLIC_API_URL}` + `/hls/${videoId}/manifest`;
  },
};
