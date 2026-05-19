package com.novix.content_service.dto.request;

import lombok.Data;

@Data
public class SearchRequest {
    private String query;
    private String category;
    private String uploaderUserId;
    private int page = 0;
    private int size = 20;
}
