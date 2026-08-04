'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { VideoGrid } from '@/components/video/VideoGrid';
import { searchService } from '@/services/search.service';
import { Video } from '@/types';
import { Search } from 'lucide-react';

export default function SearchPage() {
  const searchParams = useSearchParams();
  const query = searchParams.get('query') || '';
  const [results, setResults] = useState<Video[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [totalHits, setTotalHits] = useState(0);

  useEffect(() => {
    if (!query) return;

    const search = async () => {
      setIsLoading(true);
      try {
        const data = await searchService.search(query);
        setResults(data.results as unknown as Video[]);
        setTotalHits(data.totalHits);
      } catch (error) {
        console.error('Search failed:', error);
      } finally {
        setIsLoading(false);
      }
    };

    search();
  }, [query]);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Search className="h-5 w-5 text-gray-400" />
        <div>
          <h1 className="text-xl font-bold text-white">{query ? `Results for "${query}"` : 'Search Videos'}</h1>
          {totalHits > 0 && <p className="text-sm text-gray-400">{totalHits.toLocaleString()} results found</p>}
        </div>
      </div>

      <VideoGrid
        videos={results}
        isLoading={isLoading}
        emptyMessage={query ? `No results found for "${query}"` : 'Enter a search term'}
      />
    </div>
  );
}
