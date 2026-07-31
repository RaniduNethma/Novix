import api from '@/lib/api';
import { SearchResponse } from '@/types';

export const searchService = {
  async search(query: string, page = 0, size = 20): Promise<SearchResponse> {
    const response = await api.get<SearchResponse>(
      `/api/v1/search?query=${encodeURIComponent(query)}` + `&page=${page}&size=${size}`,
    );
    return response.data;
  },

  async searchByCategory(category: string, page = 0, size = 20): Promise<SearchResponse> {
    const response = await api.get<SearchResponse>(
      `/api/v1/search/category/${category}` + `?page=${page}&size=${size}`,
    );
    return response.data;
  },
};
