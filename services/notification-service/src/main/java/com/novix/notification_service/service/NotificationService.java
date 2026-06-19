package com.novix.notification_service.service;

import com.novix.notification_service.dto.request.SendNotificationRequest;
import com.novix.notification_service.dto.rsponse.NotificationResponse;
import com.novix.notification_service.dto.rsponse.PageResponse;

public interface NotificationService {

    NotificationResponse sendNotification(SendNotificationRequest request);

    NotificationResponse getNotificationById(Long id);

    PageResponse<NotificationResponse> getUserNotifications(String userId, int page, int size);

    PageResponse<NotificationResponse> getUnreadNotifications(String userId, int page, int size);

    Long getUnreadCount(String userId);

    void markAsRead(Long id);

    void markAllAsRead(String userId);
}
