'use client';

import { useEffect, useState } from 'react';
import { StatsCard } from '@/components/charts/StatsCard';
import { Card } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { statsService } from '@/services/stats.service';
import { videoService } from '@/services/video.service';
import { Video, DashboardStats } from '@/types';
import { Spinner } from '@/components/ui/Spinner';
import { Video as VideoIcon, Tag, TrendingUp, CheckCircle } from 'lucide-react';

function getStatusVariant(status: string) {
  switch (status) {
    case 'READY':
      return 'success';
    case 'PROCESSING':
      return 'warning';
    case 'FAILED':
      return 'danger';
    case 'UPLOADING':
      return 'info';
    default:
      return 'default';
  }
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [recentVideos, setRecentVideos] = useState<Video[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [statsData, videosData] = await Promise.all([
          statsService.getDashboardStats(),
          videoService.getVideos(0, 5),
        ]);
        setStats(statsData);
        setRecentVideos(videosData.content);
      } catch (err) {
        console.error('Failed to load dashboard:', err);
      } finally {
        setIsLoading(false);
      }
    };
    load();
  }, []);

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Spinner className="h-10 w-10" />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard title="Total Videos" value={stats?.totalVideos || 0} icon={VideoIcon} color="red" />
        <StatsCard title="Total Categories" value={stats?.totalCategories || 0} icon={Tag} color="blue" />
        <StatsCard title="Ready Videos" value={stats?.readyVideos || 0} icon={CheckCircle} color="green" />
        <StatsCard title="Processing" value={stats?.processingVideos || 0} icon={TrendingUp} color="yellow" />
      </div>

      {/* Recent Videos */}
      <Card>
        <h3 className="font-bold text-white mb-4">Recent Videos</h3>
        {recentVideos.length === 0 ? (
          <p className="text-gray-500 text-sm text-center py-8">No videos yet</p>
        ) : (
          <div className="space-y-3">
            {recentVideos.map((video) => (
              <div key={video.videoId} className="flex items-center justify-between p-3 bg-gray-800 rounded-lg">
                <div className="flex-1 min-w-0 mr-4">
                  <p className="text-sm font-medium text-white truncate">{video.title}</p>
                  <p className="text-xs text-gray-400">
                    {video.uploaderUsername} • {new Date(video.createdAt).toLocaleDateString()}
                  </p>
                </div>
                <Badge variant={getStatusVariant(video.status)}>{video.status}</Badge>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
