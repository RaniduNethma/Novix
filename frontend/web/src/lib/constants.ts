export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export const APP_NAME = process.env.NEXT_PUBLIC_APP_NAME || 'Novix';

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  WATCH: (videoId: string) => `/watch/${videoId}`,
  SEARCH: '/search',
  PROFILE: '/profile',
} as const;

export const API_ENDPOINTS = {
  // Auth
  LOGIN: '/api/v1/auth/login',
  REGISTER: '/api/v1/auth/register',
  REFRESH: '/api/v1/auth/refresh-token',
  LOGOUT: '/api/v1/auth/logout',

  // Users
  ME: '/api/v1/users/me',

  // Videos
  VIDEOS: '/api/v1/videos',
  VIDEO: (id: string) => `/api/v1/videos/${id}`,

  // Search
  SEARCH: '/api/v1/search',

  // Recommendations
  RECOMMENDATIONS: '/api/v1/recommendations/me',
  TRENDING: '/api/v1/recommendations/trending',

  // Streaming
  STREAM_START: '/api/v1/stream/start',
  STREAM_END: (sessionId: string) => `/api/v1/stream/session/${sessionId}`,
  HLS_MANIFEST: (videoId: string) => `/hls/${videoId}/manifest`,
} as const;
