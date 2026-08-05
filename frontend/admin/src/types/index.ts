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
  categories: string[];
  createdAt: string;
  publishedAt: string | null;
}

export type VideoStatus = 'UPLOADING' | 'PROCESSING' | 'READY' | 'FAILED' | 'DELETED';

export type VideoVisibility = 'PUBLIC' | 'PRIVATE' | 'UNLISTED';

export interface Category {
  id: number;
  name: string;
  description: string;
  slug: string;
  createdAt: string;
}

export interface PageResponse<T> {
  content: T[];
  page: number;
  size: number;
  totalElements: number;
  totalPages: number;
  last: boolean;
}

export interface DashboardStats {
  totalVideos: number;
  totalUsers: number;
  totalCategories: number;
  readyVideos: number;
  processingVideos: number;
  failedVideos: number;
}
