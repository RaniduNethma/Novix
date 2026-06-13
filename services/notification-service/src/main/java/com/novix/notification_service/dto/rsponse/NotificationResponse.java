package com.novix.notification_service.dto.rsponse;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class NotificationResponse {
    private Long id;
    private String userId;
    private String type;
    private String status;
    private String title;
    private String message;
    private String referenceId;
    private Boolean isRead;
    private LocalDateTime sentAt;
    private LocalDateTime createdAt;
}
