package com.novix.recommendation_service.consumer;

import com.novix.recommendation_service.event.VideoWatchedEvent;
import com.novix.recommendation_service.service.RecommendationService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Slf4j
public class WatchEventConsumer {
    private final RecommendationService recommendationService;

    @KafkaListener(
            topics = "video.watched",
            groupId = "recommendation-service"
    )

    public void handleVideoWatched(VideoWatchedEvent event) {
        log.info("Received video.watched event: user={} video={}", event.getUserId(), event.getVideoId());
        recommendationService.processWatchEvent(event);
    }
}
