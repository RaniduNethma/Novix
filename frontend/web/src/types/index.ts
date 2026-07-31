export interface User {
  id: number;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  profilePicture: string | null;
  isActive: boolean;
  isEmailVerified: boolean;
  roles: string[];
  createdAt: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  user: User;
}

export interface Video {
  id: number;
  videoId: string;
  title: string;
  description: string;
  thumbnailUrl: string | null;
  uploaderUserId: string;
  uploaderUsername: string;
  status: VideoStatus;
  visibility: VideoVisibility;
  duration: number | null;
  viewCount: number;
  likeCount: number;
  masterPlaylistPath: string | null;
  categories: string[];
  createdAt: string;
  publishedAt: string | null;
}

export type VideoStatus = 'UPLOADING' | 'PROCESSING' | 'READY' | 'FAILED' | 'DELETED';

export type VideoVisibility = 'PUBLIC' | 'PRIVATE' | 'UNLISTED';

export interface PageResponse<T> {
  content: T[];
  page: number;
  size: number;
  totalElements: number;
  totalPages: number;
  last: boolean;
}

export interface SearchResponse {
  results: Video[];
  totalHits: number;
  page: number;
  size: number;
  query: string;
}

export interface VideoRecommendation {
  videoId: string;
  title: string;
  thumbnailUrl: string | null;
  categories: string[];
  score: number;
  reason: string;
}

export interface RecommendationResponse {
  recommendations: VideoRecommendation[];
  count: number;
}

export interface StreamResponse {
  sessionId: string;
  manifestUrl: string;
  thumbnailUrl: string | null;
  expiresIn: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  password: string;
}
