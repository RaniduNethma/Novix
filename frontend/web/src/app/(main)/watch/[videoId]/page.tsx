'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { VideoPlayer } from '@/components/video/VideoPlayer';
import { Spinner } from '@/components/ui/Spinner';
import { Video, StreamResponse } from '@/types';
import { videoService } from '@/services/video.service';
import { streamService } from '@/services/stream.service';
import { useAuth } from '@/context/AuthContext';
import { Eye, ThumbsUp, Calendar } from 'lucide-react';

export default function WatchPage() {
  const { videoId } = useParams<{ videoId: string }>();
  const { isAuthenticated } = useAuth();
  const [video, setVideo] = useState<Video | null>(null);
  const [stream, setStream] = useState<StreamResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadVideo = async () => {
      setIsLoading(true);
      try {
        // Load video metadata
        const videoData = await videoService.getVideoById(videoId);
        setVideo(videoData);

        // Start stream session if authenticated
        if (isAuthenticated) {
          const streamData = await streamService.startStream(videoId);
          setStream(streamData);
        }
      } catch (err) {
        setError('Video not found or not available.');
      } finally {
        setIsLoading(false);
      }
    };

    if (videoId) loadVideo();

    // Cleanup — end stream session on unmount
    return () => {
      if (stream?.sessionId) {
        streamService.endStream(stream.sessionId).catch(console.error);
      }
    };
  }, [videoId, isAuthenticated]);

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Spinner className="h-12 w-12" />
      </div>
    );
  }

  if (error || !video) {
    return (
      <div className="text-center py-20">
        <p className="text-red-400 text-lg">{error}</p>
      </div>
    );
  }

  const manifestUrl = stream?.manifestUrl || streamService.getManifestUrl(videoId);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      {/* Player + Info */}
      <div className="lg:col-span-2 space-y-6">
        {/* Video Player */}
        <VideoPlayer manifestUrl={manifestUrl} sessionId={stream?.sessionId || ''} />

        {/* Video Info */}
        <div className="space-y-4">
          <h1 className="text-2xl font-bold text-white">{video.title}</h1>

          <div className="flex items-center justify-between flex-wrap gap-4 pb-4 border-b border-gray-800">
            <div className="flex items-center gap-4 text-sm text-gray-400">
              <span className="flex items-center gap-1">
                <Eye className="h-4 w-4" />
                {video.viewCount?.toLocaleString() || 0} views
              </span>
              <span className="flex items-center gap-1">
                <ThumbsUp className="h-4 w-4" />
                {video.likeCount?.toLocaleString() || 0}
              </span>
              {video.publishedAt && (
                <span className="flex items-center gap-1">
                  <Calendar className="h-4 w-4" />
                  {new Date(video.publishedAt).toLocaleDateString()}
                </span>
              )}
            </div>
          </div>

          {/* Uploader */}
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-red-600 flex items-center justify-center font-bold">
              {video.uploaderUsername?.[0]?.toUpperCase()}
            </div>
            <div>
              <p className="font-medium text-white">{video.uploaderUsername}</p>
              <p className="text-xs text-gray-400">Creator</p>
            </div>
          </div>

          {/* Description */}
          {video.description && (
            <div className="bg-gray-900 rounded-xl p-4">
              <p className="text-gray-300 text-sm whitespace-pre-wrap">{video.description}</p>
            </div>
          )}

          {/* Categories */}
          {video.categories?.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {video.categories.map((cat) => (
                <span key={cat} className="bg-gray-800 text-gray-300 text-xs px-3 py-1 rounded-full">
                  {cat}
                </span>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Sidebar — Up Next placeholder */}
      <div className="lg:col-span-1">
        <h3 className="font-bold text-white mb-4">Up Next</h3>
        <div className="space-y-4 text-sm text-gray-500 text-center py-8">More videos coming soon</div>
      </div>
    </div>
  );
}
