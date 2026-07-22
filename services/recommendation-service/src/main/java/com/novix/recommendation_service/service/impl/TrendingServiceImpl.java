package com.novix.recommendation_service.service.impl;

import com.novix.recommendation_service.document.TrendingVideo;
import com.novix.recommendation_service.dto.response.VideoRecommendation;
import com.novix.recommendation_service.repository.TrendingVideoRepository;
import com.novix.recommendation_service.service.TrendingService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j
public class TrendingServiceImpl implements TrendingService {
    private final TrendingVideoRepository trendingVideoRepository;

    @Override
    public void recordView(String videoId, String title, List<String> categories) {
        Optional<TrendingVideo> existing = trendingVideoRepository.findByVideoId(videoId);

        TrendingVideo trendingVideo;
        if (existing.isPresent()) {
            trendingVideo = existing.get();
            trendingVideo.setViewCount(trendingVideo.getViewCount() + 1);
            trendingVideo.setRecentViewCount(trendingVideo.getViewCount() + 1);
        } else {
            trendingVideo = TrendingVideo.builder()
                    .videoId(videoId)
                    .title(title)
                    .categories(categories)
                    .viewCount(1L)
                    .recentViewCount(1L)
                    .build();
        }

        // Simple trending score : recent view weighted higher
        trendingVideo.setTrendingScore(trendingVideo.getRecentViewCount() * 2.0 + trendingVideo.getViewCount() * 0.5);
        trendingVideo.setLastUpdated(LocalDateTime.now());
        trendingVideoRepository.save(trendingVideo);
    }

    @Override
    public List<VideoRecommendation> getTrendingVideos(int limit) {
        List<TrendingVideo> trendingVideo = trendingVideoRepository.findTop20ByOrderByTrendingScoreDesc();

        return trendingVideo.stream()
                .limit(limit)
                .map(t -> VideoRecommendation.builder()
                        .videoId(t.getVideoId())
                        .title(t.getTitle())
                        .thumbnailUrl(t.getThumbnailUrl())
                        .categories(t.getCategories())
                        .score(t.getTrendingScore())
                        .reason("trending")
                        .build())
                .toList();
    }
}
