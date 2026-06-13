package com.novix.notification_service.dto.request;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

@Data
public class SendNotificationRequest {

    @NotBlank(message = "User ID is required")
    private String userId;

    private String email;

    @NotNull(message = "Notification type is required")
    private String type;

    @NotBlank(message = "Title is required")
    private String title;

    @NotBlank(message = "Message is required")
    private String message;

    private String referenceId;
}
