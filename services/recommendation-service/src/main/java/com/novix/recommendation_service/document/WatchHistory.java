package com.novix.recommendation_service.document;

import lombok.*;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.List;

@Document(collation = "watch_history")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class WatchHistory {

    @Id
    private String id;

    @Indexed
    private String userId;

    @Indexed
    private String videoId;

    private String title;

    private Long watchDurationSeconds;

    private Long videoDurationSeconds;

    private Double completionPercentage;

    private List<String> categories;

    @Indexed
    private LocalDateTime watchedAt;
}
