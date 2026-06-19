package com.novix.notification_service.consumer;

import com.novix.notification_service.dto.request.SendNotificationRequest;
import com.novix.notification_service.event.VideoProcessedEvent;
import com.novix.notification_service.event.VideoUploadedEvent;
import com.novix.notification_service.service.NotificationService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Slf4j
public class VideoEventConsumer {
    private final NotificationService notificationService;

    @KafkaListener(
            topics = "video.uploaded",
            groupId = "notification-service"
    )
    public void handleVideoUploaded(VideoUploadedEvent event) {
        log.info("Received video.uploaded event: {}", event.getVideoId());
        SendNotificationRequest request = new SendNotificationRequest();
        request.setUserId(event.getUserId());
        request.setType("VIDEO_UPLOADED");
        request.setTitle("Video Upload Started");
        request.setMessage(String.format(
                "Your video '%s' has been uploaded and is now being processed.",
                event.getTitle() != null ? event.getTitle() : event.getVideoId()
        ));
        request.setReferenceId(event.getVideoId());
        notificationService.sendNotification(request);
    }

    @KafkaListener(
            topics = "video.processed",
            groupId = "notification-service"
    )
    public void handleVideoProcessed(VideoProcessedEvent event) {
        log.info("Received video.processed event: {}", event.getVideoId());
        boolean success = "COMPLETED".equals(event.getStatus());
        SendNotificationRequest request = new SendNotificationRequest();
        request.setUserId(event.getUserId());
        request.setEmail(event.getEmail());
        request.setType("VIDEO_PROCESSED");
        request.setTitle(success ? "Your Video is Ready!" : "Video Processing Failed");
        request.setMessage(success ? String.format("Your video '%s' has been processed successfully and is now available for streaming.",
                event.getTitle() != null ? event.getTitle() : event.getVideoId()) : String.format(
                "Unfortunately, your video '%s' failed to process. Please try uploading again.",
                event.getTitle() != null ? event.getTitle() : event.getVideoId())
        );
        request.setReferenceId(event.getVideoId());
        notificationService.sendNotification(request);
    }
}
