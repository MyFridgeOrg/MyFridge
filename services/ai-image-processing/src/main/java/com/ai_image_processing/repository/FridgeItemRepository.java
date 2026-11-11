package com.ai_image_processing.repository;

import com.ai_image_processing.models.FridgeItemEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface FridgeItemRepository extends JpaRepository<FridgeItemEntity, String> {
    Optional<FridgeItemEntity> findByNameIgnoreCase(String name);
}