import api from '@/lib/api';
import { Category, PageResponse } from '@/types';

export const categoryService = {
  async getCategories(page = 0, size = 20): Promise<PageResponse<Category>> {
    const response = await api.get<PageResponse<Category>>(`/api/v1/categories?page=${page}&size=${size}`);
    return response.data;
  },

  async createCategory(name: string, description: string, slug: string): Promise<Category> {
    const response = await api.post<Category>(
      `/api/v1/categories?name=${encodeURIComponent(name)}` +
        `&description=${encodeURIComponent(description)}` +
        `&slug=${encodeURIComponent(slug)}`,
    );
    return response.data;
  },

  async deleteCategory(id: number): Promise<void> {
    await api.delete(`/api/v1/categories/${id}`);
  },
};
