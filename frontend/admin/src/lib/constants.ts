export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export const ADMIN_ROUTES = {
  LOGIN: '/login',
  DASHBOARD: '/',
  VIDEOS: '/videos',
  USERS: '/users',
  CATEGORIES: '/categories',
  NOTIFICATIONS: '/notifications',
} as const;

export const API_ENDPOINTS = {
  // Auth
  LOGIN: '/api/v1/auth/login',
  ME: '/api/v1/users/me',

  // Videos
  VIDEOS: '/api/v1/videos',
  VIDEO: (id: string) => `/api/v1/videos/${id}`,
  VIDEO_STATUS: (id: string) => `/api/v1/videos/${id}/status`,

  // Users
  USER: (id: number) => `/api/v1/users/${id}`,

  // Categories
  CATEGORIES: '/api/v1/categories',
  CATEGORY: (id: number) => `/api/v1/categories/${id}`,

  // Notifications
  SEND_NOTIFICATION: '/api/v1/notifications/send',

  // Recommendations
  TRENDING: '/api/v1/recommendations/trending',
} as const;
