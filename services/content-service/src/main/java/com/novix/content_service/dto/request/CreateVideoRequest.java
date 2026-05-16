package com.novix.content_service.dto.request;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

import java.util.Set;

@Data
public class CreateVideoRequest {

    @NotBlank(message = "Video ID is required")
    private String videoId;

    @NotBlank(message = "Title is required")
    @Size(min = 3, max = 100, message = "Title must be between 3 and 100 characters")
    private String title;

    @Size(max = 5000, message = "Description cannot exceed 5000 characters")
    private String description;

    private String thumbnailUrl;

    private String uploaderUserId;

    private String uploaderUsername;

    private String visibility;

    private String rawVideoPath;

    private Set<Long> categoryIds;
}
