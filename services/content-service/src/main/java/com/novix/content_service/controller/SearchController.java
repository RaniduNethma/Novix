package com.novix.content_service.controller;

import com.novix.content_service.dto.request.SearchRequest;
import com.novix.content_service.dto.response.SearchResponse;
import com.novix.content_service.service.SearchService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/search")
@RequiredArgsConstructor
public class SearchController {
    private final SearchService searchService;

    @GetMapping
    public ResponseEntity<SearchResponse> search(
            @RequestParam(required = false) String query,
            @RequestParam(required = false) String category,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {

        SearchRequest request = new SearchRequest();
        request.setQuery(query);
        request.setCategory(category);
        request.setPage(page);
        request.setSize(size);

        return ResponseEntity.ok(searchService.search(request));
    }

    @GetMapping("/category/{category}")
    public ResponseEntity<SearchResponse> searchByCategory(
            @PathVariable String category,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        return ResponseEntity.ok(searchService.searchByCategory(category, page, size));
    }
}
