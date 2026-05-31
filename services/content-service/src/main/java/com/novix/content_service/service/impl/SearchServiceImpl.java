package com.novix.content_service.service.impl;

import com.novix.content_service.document.VideoDocument;
import com.novix.content_service.dto.request.SearchRequest;
import com.novix.content_service.dto.response.SearchResponse;
import com.novix.content_service.dto.response.VideoResponse;
import com.novix.content_service.repository.elasticsearch.VideoSearchRepository;
import com.novix.content_service.service.SearchService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class SearchServiceImpl implements SearchService {
    private final VideoSearchRepository videoSearchRepository;

    @Override
    public SearchResponse search(SearchRequest request) {
        Pageable pageable = PageRequest.of(request.getPage(), request.getSize());

        Page<VideoDocument> results;

        if (request.getCategory() != null && !request.getCategory().isEmpty()) {
            results = videoSearchRepository.findByCategoriesContaining(request.getCategory(), pageable);
        } else if (request.getQuery() != null && !request.getQuery().isEmpty()) {
            results = videoSearchRepository
                    .findByTitleContainingOrDescriptionContaining(
                            request.getQuery(),
                            request.getQuery(),
                            pageable
                    );
        } else {
            results = videoSearchRepository
                    .findByStatusAndVideoVisibility("READY", "PUBLIC", pageable);
        }

        List<VideoResponse> videoResponses = results.getContent()
                .stream()
                .map(this::documentToResponse)
                .toList();

        log.info("Search query: '{}' returned {} results", request.getQuery(), results.getTotalElements());

        return SearchResponse.builder()
                .results(videoResponses)
                .totalHits(results.getTotalElements())
                .page(request.getPage())
                .size(request.getSize())
                .query(request.getQuery())
                .build();
    }

    @Override
    public SearchResponse searchByCategory(String category, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);

        Page<VideoDocument> results = videoSearchRepository.findByCategoriesContaining(category, pageable);

        List<VideoResponse> videoResponses = results.getContent()
                .stream()
                .map(this::documentToResponse)
                .toList();

        return SearchResponse.builder()
                .results(videoResponses)
                .totalHits(results.getTotalElements())
                .page(page)
                .size(size)
                .query(category)
                .build();
    }

    private VideoResponse documentToResponse(VideoDocument doc) {
        return VideoResponse.builder()
                .videoId(doc.getVideoId())
                .title(doc.getTitle())
                .description(doc.getDescription())
                .thumbnailUrl(doc.getThumbnailUrl())
                .uploaderUserId(doc.getUploaderUserId())
                .uploaderUsername(doc.getUploaderUsername())
                .status(doc.getStatus())
                .videoVisibility(doc.getVideoVisibility())
                .viewCount(doc.getViewCount())
                .likeCount(doc.getLikeCount())
                .duration(doc.getDuration())
                .publishedAt(doc.getPublishedAt())
                .createdAt(doc.getCreatedAt())
                .build();
    }
}
