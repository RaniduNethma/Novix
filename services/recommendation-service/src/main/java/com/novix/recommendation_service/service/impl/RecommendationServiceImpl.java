package com.novix.recommendation_service.service.impl;

import com.novix.recommendation_service.document.UserPreference;
import com.novix.recommendation_service.document.WatchHistory;
import com.novix.recommendation_service.dto.response.RecommendationResponse;
import com.novix.recommendation_service.dto.response.VideoRecommendation;
import com.novix.recommendation_service.event.VideoWatchedEvent;
import com.novix.recommendation_service.repository.UserPreferenceRepository;
import com.novix.recommendation_service.repository.WatchHistoryRepository;
import com.novix.recommendation_service.service.RecommendationService;
import com.novix.recommendation_service.service.TrendingService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.*;

@Service
@RequiredArgsConstructor
@Slf4j
public class RecommendationServiceImpl implements RecommendationService {
    private final WatchHistoryRepository watchHistoryRepository;
    private final UserPreferenceRepository userPreferenceRepository;
    private final TrendingService trendingService;

    @Override
    public void processWatchEvent(VideoWatchedEvent event) {

        // Calculate completion percentage
        double completion = 0.0;
        if (event.getVideoDurationSeconds() != null && event.getVideoDurationSeconds() > 0) {
            completion = (event.getWatchDurationSeconds() * 100.0) / event.getVideoDurationSeconds();
        }

        // Save watch history
        WatchHistory history = WatchHistory.builder()
                .userId(event.getUserId())
                .videoId(event.getVideoId())
                .title(event.getTitle())
                .watchDurationSeconds(event.getWatchDurationSeconds())
                .videoDurationSeconds(event.getVideoDurationSeconds())
                .completionPercentage(completion)
                .categories(event.getCategories())
                .watchedAt(LocalDateTime.now())
                .build();
        watchHistoryRepository.save(history);

        // Update user category preferences
        updateUserPreferences(event.getUserId(), event.getCategories(), completion);

        // Record trending view (only count meaningful watches)
        if (completion > 30.0) {
            trendingService.recordView(event.getVideoId(), event.getTitle(), event.getCategories());
        }

        log.info("Processed watch event for user: {} video: {}", event.getUserId(), event.getVideoId());
    }

    @Override
    public RecommendationResponse getRecommendationsForUser(String userId, int limit) {
        Optional<UserPreference> preference = userPreferenceRepository.findByUserId(userId);
        List<VideoRecommendation> recommendations = new ArrayList<>();

        if (preference.isPresent() && !preference.get().getCategoryScores().isEmpty()){
            // Get top categories use likes
            List<String> topCategories = preference.get().getCategoryScores().entrySet().stream()
                    .sorted(Map.Entry.<String, Double>comparingByValue().reversed())
                    .limit(3)
                    .map(Map.Entry::getKey)
                    .toList();

            log.info("Top categories for user {}: {}", userId, topCategories);

            recommendations = trendingService.getTrendingVideos(limit).stream().map(rec -> {
                boolean matchesPreference = rec.getCategories() != null &&
                        rec.getCategories().stream().anyMatch(topCategories::contains);

                return VideoRecommendation.builder()
                        .videoId(rec.getVideoId())
                        .title(rec.getTitle())
                        .thumbnailUrl(rec.getThumbnailUrl())
                        .categories(rec.getCategories())
                        .score(rec.getScore())
                        .reason(matchesPreference ? "based on your interests" : "trending")
                        .build();
            }) .toList();
        } else {
            recommendations = trendingService.getTrendingVideos(limit);
        }

        return RecommendationResponse.builder()
                .recommendations(recommendations)
                .count(recommendations.size())
                .build();
    }

    private void updateUserPreferences(String userId, List<String> categories, double completion) {
        if (categories == null || categories.isEmpty()) return;

        UserPreference preference = userPreferenceRepository
                .findByUserId(userId)
                .orElse(UserPreference.builder()
                        .userId(userId)
                        .categoryScores(new HashMap<>())
                        .build());

        // Weight score increase by completion percentage
        double scoreIncrease = completion / 100.0;

        for (String category : categories) {
            preference.getCategoryScores().merge(category, scoreIncrease, Double::sum);
        }

        preference.setUpdatedAt(LocalDateTime.now());
        userPreferenceRepository.save(preference);
    }
}
