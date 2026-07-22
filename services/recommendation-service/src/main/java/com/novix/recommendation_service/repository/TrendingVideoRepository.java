package com.novix.recommendation_service.repository;

import com.novix.recommendation_service.document.TrendingVideo;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface TrendingVideoRepository extends MongoRepository<TrendingVideo, String> {

    Optional<TrendingVideo> findByVideoId(String videoId);

    List<TrendingVideo> findTop20ByOrderByTrendingScoreDesc();
}
