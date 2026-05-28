package com.novix.content_service.repository.jpa;

import com.novix.content_service.entity.Category;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface CategoryRepository extends JpaRepository<Category, Long> {
    Optional<Category> findBySlug(String slug);
    Optional<Category> findByName(String name);
    Boolean existsBySlug(String slug);
    Boolean existsByName(String name);
    Page<Category> findAll(Pageable pageable);
}
