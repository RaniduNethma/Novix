'use client';

import { useEffect, useState } from 'react';
import { VideoGrid } from '@/components/video/VideoGrid';
import { VideoRecommendation, Video } from '@/types';
import { recommendationService } from '@/services/recommendation.service';
import { videoService } from '@/services/video.service';
import { useAuth } from '@/context/AuthContext';
import { TrendingUp, Sparkles } from 'lucide-react';

export default function HomePage() {
  const { isAuthenticated } = useAuth();
  const [trending, setTrending] = useState<VideoRecommendation[]>([]);
  const [recommended, setRecommended] = useState<VideoRecommendation[]>([]);
  const [latest, setLatest] = useState<Video[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadContent = async () => {
      setIsLoading(true);
      try {
        // Always load trending and latest
        const [trendingData, latestData] = await Promise.all([
          recommendationService.getTrending(8),
          videoService.getPublicVideos(0, 8),
        ]);

        setTrending(trendingData.recommendations);
        setLatest(latestData.content);

        // Load personalized only if authenticated
        if (isAuthenticated) {
          const recData = await recommendationService.getMyRecommendations(8);
          setRecommended(recData.recommendations);
        }
      } catch (error) {
        console.error('Failed to load content:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadContent();
  }, [isAuthenticated]);

  return (
    <div className="space-y-12">
      {/* Hero Banner */}
      <div className="relative rounded-2xl overflow-hidden bg-gradient-to-r from-red-900/50 to-gray-900 p-8 md:p-12">
        <div className="relative z-10">
          <h1 className="text-3xl md:text-5xl font-bold text-white mb-3">Stream Anything</h1>
          <p className="text-gray-300 text-lg mb-6 max-w-xl">
            Discover and watch videos from creators around the world. High quality adaptive streaming.
          </p>
        </div>
        <div className="absolute inset-0 bg-gradient-to-r from-red-900/20 to-transparent" />
      </div>

      {/* Personalized Recommendations */}
      {isAuthenticated && recommended.length > 0 && (
        <section>
          <div className="flex items-center gap-2 mb-6">
            <Sparkles className="h-5 w-5 text-red-500" />
            <h2 className="text-xl font-bold text-white">Recommended For You</h2>
          </div>
          <VideoGrid videos={recommended} isLoading={false} />
        </section>
      )}

      {/* Trending */}
      <section>
        <div className="flex items-center gap-2 mb-6">
          <TrendingUp className="h-5 w-5 text-red-500" />
          <h2 className="text-xl font-bold text-white">Trending Now</h2>
        </div>
        <VideoGrid videos={trending} isLoading={isLoading} emptyMessage="No trending videos yet" />
      </section>

      {/* Latest Videos */}
      <section>
        <h2 className="text-xl font-bold text-white mb-6">Latest Videos</h2>
        <VideoGrid videos={latest} isLoading={isLoading} emptyMessage="No videos yet" />
      </section>
    </div>
  );
}
