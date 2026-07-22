package com.novix.recommendation_service.repository;

import com.novix.recommendation_service.document.WatchHistory;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface WatchHistoryRepository extends MongoRepository<WatchHistory, String> {

    List<WatchHistory> findByUserIdOrderByWatchedAtDesc(String userId);

    List<WatchHistory> findTop50ByUserIdOrderByWatchedAtDesc(String userId);

    boolean existsByUserIdAndVideoId(String userId, String videoId);
}
