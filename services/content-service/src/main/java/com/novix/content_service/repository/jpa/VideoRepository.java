package com.novix.content_service.repository.jpa;

import com.novix.content_service.entity.Video;
import com.novix.content_service.enums.VideoStatus;
import com.novix.content_service.enums.VideoVisibility;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface VideoRepository extends JpaRepository<Video, Long> {
    Optional<Video> findByVideoId(String videoId);
    Boolean existsByVideoId(String videoId);
    Page<Video> findByUploaderUserId (String uploaderUserId, Pageable pageable);
    Page<Video> findByStatusAndVisibility(VideoStatus status, VideoVisibility videoVisibility, Pageable pageable);
    Page<Video> findByCategoriesId(Long categoryId, Pageable pageable);
    void deleteByVideoId(String videoId);
}
