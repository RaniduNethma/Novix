package com.novix.recommendation_service.service;

import com.novix.recommendation_service.dto.response.RecommendationResponse;
import com.novix.recommendation_service.event.VideoWatchedEvent;

public interface RecommendationService {

    void processWatchEvent(VideoWatchedEvent event);

    RecommendationResponse getRecommendationsForUser(String userId, int limit);
}
