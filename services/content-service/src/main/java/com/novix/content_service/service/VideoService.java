package com.novix.content_service.service;

import com.novix.content_service.dto.request.CreateVideoRequest;
import com.novix.content_service.dto.request.UpdateVideoRequest;
import com.novix.content_service.dto.response.PageResponse;
import com.novix.content_service.dto.response.VideoResponse;
import com.novix.content_service.enums.VideoStatus;

public interface VideoService {
    VideoResponse createVideo(CreateVideoRequest request);

    VideoResponse getVideoByVideoId(String videoId);

    VideoResponse updateVideo(String videoId, UpdateVideoRequest request);

    void deleteVideo(String videoId);

    void updateVideoStatus(String videoId, VideoStatus status);

    PageResponse<VideoResponse> getAllPublicVideos(int page, int size);

    PageResponse<VideoResponse> getVideosByUploader(String uploaderUserId, int page, int size);

    PageResponse<VideoResponse> getVideosByCategory(Long categoryId, int page, int size);
}
