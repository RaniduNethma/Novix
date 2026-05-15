package com.novix.content_service.entity;

import com.novix.content_service.enums.VideoStatus;
import com.novix.content_service.enums.VideoVisibility;
import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import java.time.LocalDateTime;
import java.util.HashSet;
import java.util.Set;

@Entity
@Table(name = "videos")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class Video {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, unique = true)
    private String videoId; // UUID from processing service

    @Column(nullable = false)
    private String title;

    @Column(columnDefinition = "TEXT")
    private String description;

    private String thumbnailUrl;

    private String uploaderUserId;

    private String uploaderUserName;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    @Builder.Default
    private VideoStatus status = VideoStatus.UPLOADING;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    @Builder.Default
    private VideoVisibility videoVisibility = VideoVisibility.PUBLIC;

    private Long duration;

    private Long viewCount;

    @Builder.Default
    private Long likeCount = 0L;

    private String masterPlaylistPath;

    private String rawVideoPath;

    @CreationTimestamp
    @Column(updatable = false)
    private LocalDateTime createdAt;

    @CreationTimestamp
    private LocalDateTime updatedAt;

    @CreationTimestamp
    private LocalDateTime publishedAt;

    @ManyToMany(fetch = FetchType.LAZY)
    @JoinTable(
            name = "video_categories",
            joinColumns = @JoinColumn(name = "video_id"),
            inverseJoinColumns = @JoinColumn(name = "category_id")
    )
    @Builder.Default
    private Set<Category> categories = new HashSet<>();
}
