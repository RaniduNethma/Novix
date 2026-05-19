package com.novix.content_service.mapper;

import com.novix.content_service.document.VideoDocument;
import com.novix.content_service.dto.response.VideoResponse;
import com.novix.content_service.entity.Category;
import com.novix.content_service.entity.Video;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;

import java.util.Set;
import java.util.stream.Collectors;

@Mapper(componentModel = "spring")
public interface VideoMapper {
    @Mapping(target = "status", source = "status")
    @Mapping(target = "visibility", source = "visibility")
    @Mapping(target = "categories", source = "categories", qualifiedByName = "categoriesToStrings")
    VideoResponse toVideoResponse(Video video);

    @Mapping(target = "status", expression = "java(video.getStatus().name())")
    @Mapping(target = "visibility", expression = "java(video.getVisibility().name())")
    @Mapping(target = "categories", source = "categories", qualifiedByName = "categoriesToStringList")
    VideoDocument toVideoDocument(Video video);

    @Named("categoriesToStrings")
    default Set<String> categoriesToStrings(Set<Category> categories) {
        if (categories == null) return Set.of();
        return categories.stream().map(Category::getName).collect(Collectors.toSet());
    }

    @Named("categoriesToStringList")
    default java.util.List<String> categoriesToStringList(Set<Category> categories) {
        if (categories == null) return java.util.List.of();
        return categories.stream().map(Category::getName).collect(Collectors.toList());
    }
}
