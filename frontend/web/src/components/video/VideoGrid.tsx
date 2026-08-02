import { Video, VideoRecommendation } from '@/types';
import { VideoCard } from './VideoCard';
import { Spinner } from '@/components/ui/Spinner';

interface VideoGridProps {
  videos: (Video | VideoRecommendation)[];
  isLoading?: boolean;
  title?: string;
  emptyMessage?: string;
}

export function VideoGrid({ videos, isLoading = false, title, emptyMessage = 'No videos found' }: VideoGridProps) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Spinner className="h-10 w-10" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {title && <h2 className="text-xl font-bold text-white">{title}</h2>}
      {videos.length === 0 ? (
        <div className="text-center py-20 text-gray-500">{emptyMessage}</div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {videos.map((video) => (
            <VideoCard key={video.videoId} video={video} />
          ))}
        </div>
      )}
    </div>
  );
}
