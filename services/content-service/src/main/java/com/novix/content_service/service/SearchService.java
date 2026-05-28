package com.novix.content_service.service;

import com.novix.content_service.dto.request.SearchRequest;
import com.novix.content_service.dto.response.SearchResponse;

public interface SearchService {
    SearchResponse search(SearchRequest request);
    SearchResponse searchByCategory(String category, int page, int size);
}
