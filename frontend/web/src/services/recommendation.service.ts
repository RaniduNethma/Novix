import api from '@/lib/api';
import { RecommendationResponse } from '@/types';

export const recommendationService = {
  async getMyRecommendations(limit = 20): Promise<RecommendationResponse> {
    const response = await api.get<RecommendationResponse>(`/api/v1/recommendations/me?limit=${limit}`);
    return response.data;
  },

  async getTrending(limit = 20): Promise<RecommendationResponse> {
    const response = await api.get<RecommendationResponse>(`/api/v1/recommendations/trending?limit=${limit}`);
    return response.data;
  },
};
