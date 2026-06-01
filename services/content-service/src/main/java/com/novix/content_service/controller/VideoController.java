package com.novix.content_service.controller;

import com.novix.content_service.dto.request.CreateVideoRequest;
import com.novix.content_service.dto.request.UpdateVideoRequest;
import com.novix.content_service.dto.response.PageResponse;
import com.novix.content_service.dto.response.VideoResponse;
import com.novix.content_service.enums.VideoStatus;
import com.novix.content_service.service.VideoService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/videos")
@RequiredArgsConstructor
public class VideoController {
    private final VideoService videoService;

    @PostMapping
    public ResponseEntity<VideoResponse> createVideo(
            @Valid @RequestBody CreateVideoRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED).body(videoService.createVideo(request));
    }

    @GetMapping("/{videoId}")
    public ResponseEntity<VideoResponse> getVideoById(
            @PathVariable String videoId) {
        return ResponseEntity.ok(videoService.getVideoByVideoId(videoId));
    }

    @PutMapping("/{videoId}")
    public ResponseEntity<VideoResponse> updateVideo(
            @PathVariable String videoId,
            @Valid @RequestBody UpdateVideoRequest request) {
        return ResponseEntity.ok(videoService.updateVideo(videoId, request));
    }

    @DeleteMapping("/{videoId}")
    @PreAuthorize("hasRole('ROLE_ADMIN')")
    public ResponseEntity<String> deleteVideo(
            @PathVariable String videoId) {
        videoService.deleteVideo(videoId);
        return ResponseEntity.ok("Video deleted successfully");
    }

    @PatchMapping("/{videoId}/status")
    public ResponseEntity<String> updateStatus(
            @PathVariable String videoId,
            @RequestParam VideoStatus status) {
        videoService.updateVideoStatus(videoId, status);
        return ResponseEntity.ok("Status updated to: " + status);
    }

    @GetMapping
    public ResponseEntity<PageResponse<VideoResponse>> getAllPublicVideos(
            @PathVariable String uploaderUserId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        return ResponseEntity.ok(videoService.getVideosByUploader(uploaderUserId, page, size));
    }

    @GetMapping("/uploader/{uploaderUserId}")
    public ResponseEntity<PageResponse<VideoResponse>> getVideosByUploader(
            @PathVariable String uploaderUserId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        return ResponseEntity.ok(videoService.getVideosByUploader(uploaderUserId, page, size));
    }

    @GetMapping("/category/{categoryId}")
    public ResponseEntity<PageResponse<VideoResponse>> getVideosByCategory(
            @PathVariable Long categoryId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        return ResponseEntity.ok(videoService.getVideosByCategory(categoryId, page, size));
    }
}
