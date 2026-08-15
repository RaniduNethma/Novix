'use client';

import { useEffect, useState } from 'react';
import { Video, VideoStatus } from '@/types';
import { videoService } from '@/services/video.service';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { Spinner } from '@/components/ui/Spinner';
import { Trash2, RefreshCw, ChevronLeft, ChevronRight } from 'lucide-react';
import toast from 'react-hot-toast';

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

export default function VideosPage() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [totalElements, setTotalElements] = useState(0);

  const loadVideos = async (p = 0) => {
    setIsLoading(true);
    try {
      const data = await videoService.getVideos(p, 10);
      setVideos(data.content);
      setTotalPages(data.totalPages);
      setTotalElements(data.totalElements);
    } catch {
      toast.error('Failed to load videos');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadVideos(page);
  }, [page]);

  const handleDelete = async (videoId: string) => {
    if (!confirm('Are you sure you want to delete this video?')) return;
    try {
      await videoService.deleteVideo(videoId);
      toast.success('Video deleted');
      loadVideos(page);
    } catch {
      toast.error('Failed to delete video');
    }
  };

  const handleStatusChange = async (videoId: string, status: VideoStatus) => {
    try {
      await videoService.updateVideoStatus(videoId, status);
      toast.success(`Status updated to ${status}`);
      loadVideos(page);
    } catch {
      toast.error('Failed to update status');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <p className="text-sm text-gray-400">{totalElements.toLocaleString()} total videos</p>
        <Button variant="secondary" size="sm" onClick={() => loadVideos(page)}>
          <RefreshCw className="h-4 w-4 mr-1" />
          Refresh
        </Button>
      </div>

      <Card className="p-0 overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center py-16">
            <Spinner className="h-8 w-8" />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-800">
                  <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                    Video
                  </th>
                  <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                    Uploader
                  </th>
                  <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                    Views
                  </th>
                  <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                    Date
                  </th>
                  <th className="text-right text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-800">
                {videos.map((video) => (
                  <tr key={video.videoId} className="hover:bg-gray-800/50 transition-colors">
                    <td className="px-6 py-4">
                      <p className="text-sm font-medium text-white truncate max-w-xs">{video.title}</p>
                      <p className="text-xs text-gray-500 mt-0.5">{video.videoId}</p>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-300">{video.uploaderUsername}</td>
                    <td className="px-6 py-4">
                      <Badge variant={getStatusVariant(video.status)}>{video.status}</Badge>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-300">{video.viewCount?.toLocaleString() || 0}</td>
                    <td className="px-6 py-4 text-sm text-gray-400">
                      {new Date(video.createdAt).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-2">
                        {video.status === 'FAILED' && (
                          <Button
                            size="sm"
                            variant="secondary"
                            onClick={() => handleStatusChange(video.videoId, 'PROCESSING')}
                          >
                            Retry
                          </Button>
                        )}
                        {video.status === 'READY' && (
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleStatusChange(video.videoId, 'DELETED')}
                          >
                            Unpublish
                          </Button>
                        )}
                        <Button size="sm" variant="danger" onClick={() => handleDelete(video.videoId)}>
                          <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-6 py-4 border-t border-gray-800">
            <p className="text-sm text-gray-400">
              Page {page + 1} of {totalPages}
            </p>
            <div className="flex gap-2">
              <Button size="sm" variant="secondary" onClick={() => setPage((p) => p - 1)} disabled={page === 0}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <Button
                size="sm"
                variant="secondary"
                onClick={() => setPage((p) => p + 1)}
                disabled={page >= totalPages - 1}
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}
