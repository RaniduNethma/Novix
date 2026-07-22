package com.novix.recommendation_service.controller;

import com.novix.recommendation_service.dto.response.RecommendationResponse;
import com.novix.recommendation_service.dto.response.VideoRecommendation;
import com.novix.recommendation_service.service.RecommendationService;
import com.novix.recommendation_service.service.TrendingService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("api/v1/recommendations")
@RequiredArgsConstructor
public class RecommendationController {
    private final RecommendationService recommendationService;
    private final TrendingService trendingService;

    @GetMapping("/me")
    public ResponseEntity<RecommendationResponse> getMyRecommendations(
            @AuthenticationPrincipal UserDetails userDetails,
            @RequestParam(defaultValue = "20") int limit) {
        return ResponseEntity.ok(recommendationService.getRecommendationsForUser(userDetails.getUsername(), limit));
    }

    @GetMapping("/trending")
    public ResponseEntity<List<VideoRecommendation>> getTrending(
            @RequestParam(defaultValue = "20") int limit) {
        return ResponseEntity.ok(trendingService.getTrendingVideos(limit));
    }
}
