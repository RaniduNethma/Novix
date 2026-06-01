package com.novix.content_service.service;

import com.novix.content_service.dto.request.CreateCategoryRequest;
import com.novix.content_service.dto.response.PageResponse;
import com.novix.content_service.entity.Category;
import java.util.List;

public interface CategoryService {
    Category createCategory(CreateCategoryRequest request);

    Category getCategoryById(Long id);

    Category getCategoryBySlug(String slug);

    PageResponse<Category> getAllCategories(int page, int size);

    Category updateCategory(Long id, String name, String description, String slug);

    void deleteCategory(Long id);
}
