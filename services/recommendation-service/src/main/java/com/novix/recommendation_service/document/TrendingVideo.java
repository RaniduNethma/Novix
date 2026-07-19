package com.novix.recommendation_service.document;

import lombok.*;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.List;

@Document(collation = "trending_videos")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class TrendingVideo {

    @Id
    private String id;

    @Indexed(unique = true)
    private String videoId;

    private String title;

    private String thumbnailUrl;

    private List<String> categories;

    @Builder.Default
    private Long viewCount = 0L;

    @Builder.Default
    private Long recentViewCount = 0L;

    private Double trendingScore;

    private LocalDateTime lastUpdated;
}
