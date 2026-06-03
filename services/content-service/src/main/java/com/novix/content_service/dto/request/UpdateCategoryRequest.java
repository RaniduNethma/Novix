package com.novix.content_service.dto.request;

import lombok.Data;

@Data
public class UpdateCategoryRequest {
    private Long id;
    private String name;
    private String description;
    private String slug;
}
