package com.novix.recommendation_service.service;

import com.novix.recommendation_service.dto.response.VideoRecommendation;

import java.util.List;

public interface TrendingService {

    void recordView(String videoId, String title, List<String> categories);

    List<VideoRecommendation> getTrendingVideos(int limit);
}
