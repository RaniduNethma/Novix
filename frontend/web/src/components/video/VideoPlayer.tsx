'use client';

import { useEffect, useRef, useState } from 'react';
import Hls from 'hls.js';
import { Spinner } from '@/components/ui/Spinner';

interface VideoPlayerProps {
  manifestUrl: string;
  sessionId: string;
  onEnded?: () => void;
}

export function VideoPlayer({ manifestUrl, sessionId, onEnded }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video || !manifestUrl) return;

    setIsLoading(true);
    setError(null);

    if (Hls.isSupported()) {
      const hls = new Hls({
        enableWorker: true,
        lowLatencyMode: false,
      });

      hlsRef.current = hls;
      hls.loadSource(manifestUrl);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        setIsLoading(false);
        video.play().catch(() => {});
      });

      hls.on(Hls.Events.ERROR, (_, data) => {
        if (data.fatal) {
          setError('Failed to load video stream.');
          setIsLoading(false);
        }
      });
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      // Safari native HLS support
      video.src = manifestUrl;
      video.addEventListener('loadedmetadata', () => {
        setIsLoading(false);
        video.play().catch(() => {});
      });
    } else {
      setError('HLS is not supported in this browser.');
      setIsLoading(false);
    }

    return () => {
      hlsRef.current?.destroy();
    };
  }, [manifestUrl]);

  return (
    <div className="relative w-full aspect-video bg-black rounded-xl overflow-hidden">
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center">
          <Spinner className="h-12 w-12" />
        </div>
      )}

      {error && (
        <div className="absolute inset-0 flex items-center justify-center">
          <p className="text-red-400 text-sm">{error}</p>
        </div>
      )}

      <video ref={videoRef} className="w-full h-full" controls onEnded={onEnded} playsInline />
    </div>
  );
}
