package com.novix.recommendation_service.event;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class VideoWatchedEvent {
    private String userId;
    private String videoId;
    private String title;
    private Long watchDurationSeconds;
    private Long videoDurationSeconds;
    private List<String> categories;
}
