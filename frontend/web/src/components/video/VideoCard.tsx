import Link from 'next/link';
import { Video, VideoRecommendation } from '@/types';
import { Play, Clock } from 'lucide-react';
import { ROUTES } from '@/lib/constants';

interface VideoCardProps {
  video: Video | VideoRecommendation;
}

function formatDuration(seconds: number | null): string {
  if (!seconds) return '';
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  }
  return `${m}:${s.toString().padStart(2, '0')}`;
}

function isVideo(v: Video | VideoRecommendation): v is Video {
  return 'uploaderUsername' in v;
}

export function VideoCard({ video }: VideoCardProps) {
  const videoId = video.videoId;

  return (
    <Link href={ROUTES.WATCH(videoId)} className="group block">
      <div className="relative aspect-video bg-gray-800 rounded-xl overflow-hidden mb-3">
        {video.thumbnailUrl ? (
          <img
            src={video.thumbnailUrl}
            alt={video.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center bg-gray-800">
            <Play className="h-12 w-12 text-gray-600" />
          </div>
        )}

        {/* Play overlay */}
        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-colors flex items-center justify-center">
          <div className="opacity-0 group-hover:opacity-100 transition-opacity bg-white/20 backdrop-blur-sm rounded-full p-3">
            <Play className="h-6 w-6 text-white fill-white" />
          </div>
        </div>

        {/* Duration badge */}
        {isVideo(video) && video.duration && (
          <div className="absolute bottom-2 right-2 bg-black/80 text-white text-xs px-1.5 py-0.5 rounded flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {formatDuration(video.duration)}
          </div>
        )}
      </div>

      {/* Video Info */}
      <div className="space-y-1">
        <h3 className="font-medium text-white line-clamp-2 group-hover:text-red-400 transition-colors">
          {video.title}
        </h3>

        {isVideo(video) && <p className="text-sm text-gray-400">{video.uploaderUsername}</p>}

        <div className="flex items-center gap-2 text-xs text-gray-500">
          {isVideo(video) && video.viewCount > 0 && <span>{video.viewCount.toLocaleString()} views</span>}
          {video.categories?.length > 0 && (
            <span className="bg-gray-800 px-2 py-0.5 rounded-full">{video.categories[0]}</span>
          )}
        </div>
      </div>
    </Link>
  );
}
