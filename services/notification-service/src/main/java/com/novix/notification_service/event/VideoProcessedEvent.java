package com.novix.notification_service.event;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class VideoProcessedEvent {
    private String videoId;
    private String userId;
    private String email;
    private String title;
    private String outputPath;
    private String status;
}
