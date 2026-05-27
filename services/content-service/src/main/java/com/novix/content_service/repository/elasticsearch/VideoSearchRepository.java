package com.novix.content_service.repository.elasticsearch;

import com.novix.content_service.document.VideoDocument;
import com.novix.content_service.entity.Video;
import com.novix.content_service.enums.VideoStatus;
import com.novix.content_service.enums.VideoVisibility;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.elasticsearch.repository.ElasticsearchRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoSearchRepository extends ElasticsearchRepository<VideoDocument, String> {
    Page<VideoDocument> findByTitleContainingOrDescriptionContaining(String title, String description, Pageable pageable);

    Page<VideoDocument> findByStatusAndVideoVisibility(VideoStatus status, VideoVisibility videoVisibility, Pageable pageable);

    Page<VideoDocument> findByCategoriesContaining(String category, Pageable pageable);

    Page<VideoDocument> findByUploaderUserId(String uploaderUserId, Pageable pageable);
}
