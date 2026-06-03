package com.novix.content_service.service.impl;

import com.novix.content_service.dto.request.CreateCategoryRequest;
import com.novix.content_service.dto.request.UpdateCategoryRequest;
import com.novix.content_service.dto.response.PageResponse;
import com.novix.content_service.entity.Category;
import com.novix.content_service.exception.CategoryNotFoundException;
import com.novix.content_service.repository.jpa.CategoryRepository;
import com.novix.content_service.service.CategoryService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
@Slf4j
public class CategoryServiceImpl implements CategoryService {
    private final CategoryRepository categoryRepository;

    @Override
    @Transactional
    public Category createCategory(CreateCategoryRequest request) {
        if (categoryRepository.existsByName(request.getName())){
            throw new RuntimeException("Category " + request.getName() + " already exists.");
        }

        Category category = Category.builder()
                .name(request.getName())
                .description(request.getDescription())
                .slug(request.getSlug())
                .build();
        Category saved = categoryRepository.save(category);

        log.info("Category created: {}", saved.getName());
        return saved;
    }

    @Override
    public Category getCategoryById(Long id) {
        return categoryRepository.findById(id)
                .orElseThrow(() -> new CategoryNotFoundException("Category not found with id: " + id));
    }

    @Override
    public Category getCategoryBySlug(String slug) {
        return categoryRepository.findBySlug(slug)
                .orElseThrow(() -> new CategoryNotFoundException("Category not found with slug: " + slug));
    }

    @Override
    public PageResponse<Category> getAllCategories(int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("name").ascending());

        Page<Category> categoryPage = categoryRepository.findAll(pageable);

        return PageResponse.<Category>builder()
                .content(categoryPage.getContent())
                .page(page)
                .size(size)
                .totalElements(categoryPage.getTotalElements())
                .totalPages(categoryPage.getTotalPages())
                .last(categoryPage.isLast())
                .build();
    }

    @Override
    @Transactional
    public Category updateCategory(UpdateCategoryRequest request) {
        Category category = getCategoryById(request.getId());

        if (request.getName() != null) category.setName(request.getName());
        if (request.getDescription() != null) category.setDescription(request.getDescription());
        if (request.getSlug() != null) category.setSlug(request.getSlug());

        Category saved = categoryRepository.save(category);

        log.info("Category updated: {}" + saved.getName());
        return saved;
    }

    @Override
    @Transactional
    public void deleteCategory(Long id) {
        Category category = getCategoryById(id);
        categoryRepository.delete(category);
        log.info("Category deleted: {}", category.getName());
    }
}
