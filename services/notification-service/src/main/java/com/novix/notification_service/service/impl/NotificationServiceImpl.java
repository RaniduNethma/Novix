package com.novix.notification_service.service.impl;

import com.novix.notification_service.dto.request.SendNotificationRequest;
import com.novix.notification_service.dto.rsponse.NotificationResponse;
import com.novix.notification_service.dto.rsponse.PageResponse;
import com.novix.notification_service.entity.Notification;
import com.novix.notification_service.enums.NotificationStatus;
import com.novix.notification_service.enums.NotificationType;
import com.novix.notification_service.exception.NotificationNotFoundException;
import com.novix.notification_service.mapper.NotificationMapper;
import com.novix.notification_service.repository.NotificationRepository;
import com.novix.notification_service.service.EmailService;
import com.novix.notification_service.service.NotificationService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

@Service
@RequiredArgsConstructor
@Slf4j
public class NotificationServiceImpl implements NotificationService {
    private final NotificationRepository notificationRepository;
    private final NotificationMapper notificationMapper;
    private final EmailService emailService;

    @Override
    @Transactional
    public NotificationResponse sendNotification(SendNotificationRequest request) {
        NotificationType type = NotificationType.valueOf(request.getType().toUpperCase());

        Notification notification = Notification.builder()
                .userId(request.getUserId())
                .email(request.getEmail())
                .type(type)
                .title(request.getTitle())
                .message(request.getMessage())
                .referencedId(request.getReferenceId())
                .status(NotificationStatus.PENDING)
                .build();

        Notification saved = notificationRepository.save(notification);

        // Send email if email provided
        if (request.getEmail() != null && !request.getEmail().isEmpty()) {
            try {
                emailService.sendEmail(request.getEmail(), request.getTitle(), request.getMessage());
                saved.setStatus(NotificationStatus.SENT);
                saved.setSentAt(LocalDateTime.now());
            } catch (Exception e) {
                saved.setStatus(NotificationStatus.FAILED);
                saved.setErrorMessage(e.getMessage());
                log.error("Failed to send notification email: {}", e.getMessage());
            }
            notificationRepository.save(saved);
        }

        log.info("Notification sent to user: {}", request.getUserId());
        return notificationMapper.toResponse(saved);
    }

    @Override
    public NotificationResponse getNotificationById(Long id) {
        Notification notification = notificationRepository.findById(id)
                .orElseThrow(() -> new NotificationNotFoundException("Notification not found: " + id));
        return notificationMapper.toResponse(notification);
    }

    @Override
    public PageResponse<NotificationResponse> getUserNotifications(String userId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("createdAt").descending());
        Page<Notification> notificationPage = notificationRepository.findByUserId(userId, pageable);
        return buildPageResponse(notificationPage, page, size);
    }

    @Override
    public PageResponse<NotificationResponse> getUnreadNotifications(String userId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("createdAt").descending());
        Page<Notification> notificationPage = notificationRepository.findByUserIdAndIsRead(userId, false, pageable);
        return buildPageResponse(notificationPage, page, size);
    }

    @Override
    public Long getUnreadCount(String userId) {
        return notificationRepository.countByUserIdAndIsRead(userId, false);
    }

    @Override
    @Transactional
    public void markAsRead(Long id) {
        Notification notification = notificationRepository.findById(id)
                .orElseThrow(() -> new NotificationNotFoundException("Notification not found: " + id));
        notification.setRead(true);
        notification.setStatus(NotificationStatus.READ);
        notificationRepository.save(notification);
    }

    @Override
    @Transactional
    public void markAllAsRead(String userId) {
        notificationRepository.markAllAsReadByUserId(userId);
        log.info("Marked all notifications as read for user: {}", userId);
    }

    private PageResponse<NotificationResponse> buildPageResponse(Page<Notification> page, int pageNum, int size) {
        return PageResponse.<NotificationResponse>builder()
                .content(page.getContent().stream().map(notificationMapper::toResponse).toList())
                .page(pageNum)
                .size(size)
                .totalElements(page.getTotalElements())
                .totalPages(page.getTotalPages())
                .last(page.isLast())
                .build();
    }
}
