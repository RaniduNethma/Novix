package com.novix.content_service.repository.elasticsearch;

import com.novix.content_service.document.VideoDocument;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.elasticsearch.repository.ElasticsearchRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoSearchRepository extends ElasticsearchRepository<VideoDocument, String> {
    Page<VideoDocument> findByTitleContainingOrDescriptionContaining(String title, String description, Pageable pageable);

    Page<VideoDocument> findByStatusAndVisibility(String status, String visibility, Pageable pageable);

    Page<VideoDocument> findByCategoriesContaining(String category, Pageable pageable);

    Page<VideoDocument> findByUploaderUserId(String uploaderUserId, Pageable pageable);
}
