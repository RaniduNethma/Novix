package com.novix.content_service.dto.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.Set;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class VideoResponse {
    private Long id;
    private String videoId;
    private String title;
    private String description;
    private String thumbnailUrl;
    private String uploaderUserId;
    private String uploaderUsername;
    private String status;
    private String videoVisibility;
    private Long duration;
    private Long viewCount;
    private Long likeCount;
    private String masterPlaylistPath;
    private Set<String> categories;
    private LocalDateTime createdAt;
    private LocalDateTime publishedAt;
}
