package com.novix.recommendation_service.dto.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class VideoRecommendation {
    private String videoId;
    private String title;
    private String thumbnailUrl;
    private List<String> categories;
    private Double score;
    private String reason;
}
