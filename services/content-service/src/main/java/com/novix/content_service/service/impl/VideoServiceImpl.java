package com.novix.content_service.service.impl;

import com.novix.content_service.document.VideoDocument;
import com.novix.content_service.dto.request.CreateVideoRequest;
import com.novix.content_service.dto.request.UpdateVideoRequest;
import com.novix.content_service.dto.response.PageResponse;
import com.novix.content_service.dto.response.VideoResponse;
import com.novix.content_service.entity.Category;
import com.novix.content_service.entity.Video;
import com.novix.content_service.enums.VideoStatus;
import com.novix.content_service.enums.VideoVisibility;
import com.novix.content_service.exception.VideoNotFoundException;
import com.novix.content_service.mapper.VideoMapper;
import com.novix.content_service.repository.elasticsearch.VideoSearchRepository;
import com.novix.content_service.repository.jpa.CategoryRepository;
import com.novix.content_service.repository.jpa.VideoRepository;
import com.novix.content_service.service.VideoService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.HashSet;
import java.util.Set;

@Service
@RequiredArgsConstructor
@Slf4j
public class VideoServiceImpl implements VideoService {
    private final VideoRepository videoRepository;
    private final CategoryRepository categoryRepository;
    private final VideoSearchRepository videoSearchRepository;
    private final VideoMapper videoMapper;

    @Override
    @Transactional
    public VideoResponse createVideo(CreateVideoRequest request) {

        Set<Category> categories = new HashSet<>();
        if (request.getCategoryIds() != null) {
            categories = new HashSet<>(categoryRepository.findAllById(request.getCategoryIds()));
        }

        VideoVisibility videoVisibility = VideoVisibility.PUBLIC;
        if (request.getVisibility() != null){
            videoVisibility = VideoVisibility.valueOf(request.getVisibility().toUpperCase());
        }

        Video video = Video.builder()
                .videoId(request.getVideoId())
                .title(request.getTitle())
                .description(request.getDescription())
                .thumbnailUrl(request.getThumbnailUrl())
                .uploaderUserId(request.getUploaderUserId())
                .uploaderUserName(request.getUploaderUsername())
                .status(VideoStatus.UPLOADING)
                .videoVisibility(videoVisibility)
                .rawVideoPath(request.getRawVideoPath())
                .categories(categories)
                .viewCount(0L)
                .likeCount(0L)
                .build();
        Video saved = videoRepository.save(video);

        indexVideoDocument(saved);

        log.info("Video created: {}", saved.getVideoId());
        return videoMapper.toVideoResponse(saved);
    }

    @Override
    public VideoResponse getVideoByVideoId(String videoId) {
        Video video = videoRepository.findByVideoId(videoId)
                .orElseThrow(() -> new VideoNotFoundException("Video not found: " + videoId));
        return videoMapper.toVideoResponse(video);
    }

    @Override
    @Transactional
    public VideoResponse updateVideo(String videoId, UpdateVideoRequest request) {

        Video video = videoRepository.findByVideoId(videoId)
                .orElseThrow(() -> new VideoNotFoundException("Video not found: " + videoId));

        if (request.getTitle() != null) {
            video.setTitle(request.getTitle());
        }
        if (request.getDescription() != null) {
            video.setDescription(request.getDescription());
        }
        if (request.getThumbnailUrl() != null) {
            video.setThumbnailUrl(request.getThumbnailUrl());
        }
        if (request.getVisibility() != null) {
            video.setVideoVisibility(VideoVisibility.valueOf(request.getVisibility().toUpperCase()));
        }
        if (request.getCategoryIds() != null) {
            Set<Category> categories = new HashSet<>(categoryRepository.findAllById(request.getCategoryIds()));
            video.setCategories(categories);
        }

        Video updated = videoRepository.save(video);

        indexVideoDocument(updated);

        log.info("Video updated: {}", videoId);
        return videoMapper.toVideoResponse(updated);
    }

    @Override
    @Transactional
    public void deleteVideo(String videoId) {
        Video video = videoRepository.findByVideoId(videoId)
                .orElseThrow(() -> new VideoNotFoundException("Video not found: " + videoId));

        videoRepository.delete(video);

        videoSearchRepository.deleteById(videoId);

        log.info("Video deleted: {}", videoId);
    }

    @Override
    @Transactional
    public void updateVideoStatus(String videoId, VideoStatus status) {
        Video video = videoRepository.findByVideoId(videoId)
                .orElseThrow(() -> new VideoNotFoundException("Video not found: " + videoId));

        video.setStatus(status);

        if (status == VideoStatus.READY && video.getPublishedAt() == null) {
            video.setPublishedAt(LocalDateTime.now());
        }

        Video updated = videoRepository.save(video);

        indexVideoDocument(updated);
        log.info("Video status updated: {} -> {}", videoId, status);
    }

    @Override
    public PageResponse<VideoResponse> getAllPublicVideos(int page, int size) {
        Pageable pageable = PageRequest.of(page, size,
                Sort.by("createdAt").descending());

        Page<Video> videoPage = videoRepository
                .findByStatusAndVideoVisibility(
                        VideoStatus.READY,
                        VideoVisibility.PUBLIC,
                        pageable
                );

        return buildPageResponse(videoPage, page, size);
    }

    @Override
    public PageResponse<VideoResponse> getVideosByUploader(String uploaderUserId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("createdAt").descending());

        Page<Video> videoPage = videoRepository.findByUploaderUserId(uploaderUserId, pageable);

        return buildPageResponse(videoPage, page, size);
    }

    @Override
    public PageResponse<VideoResponse> getVideosByCategory(Long categoryId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("createdAt").descending());

        Page<Video> videoPage = videoRepository.findByCategoriesId(categoryId, pageable);

        return buildPageResponse(videoPage, page, size);
    }

    private void indexVideoDocument(Video video) {
        try {
            VideoDocument document = videoMapper.toVideoDocument(video);
            videoSearchRepository.save(document);
        } catch (Exception e) {
            log.error("Failed to index video in Elasticsearch: {}",
                    video.getVideoId(), e);
        }
    }

    private PageResponse<VideoResponse> buildPageResponse(Page<Video> videoPage, int page, int size) {
        return PageResponse.<VideoResponse>builder()
                .content(videoPage.getContent().stream().map(videoMapper::toVideoResponse).toList())
                .page(page)
                .size(size)
                .totalElements(videoPage.getTotalElements())
                .totalPages(videoPage.getTotalPages())
                .last(videoPage.isLast())
                .build();
    }
}
