package com.ai_image_processing.services;

import com.ai_image_processing.dto.FridgeItemListDTO;
import com.ai_image_processing.mapper.FridgeItemMapper;
import com.ai_image_processing.models.FridgeItemEntity;
import com.ai_image_processing.repository.FridgeItemRepository;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class FridgeItemService {

    private final FridgeItemRepository fridgeItemRepository;

    public FridgeItemService(FridgeItemRepository fridgeItemRepository){
        this.fridgeItemRepository = fridgeItemRepository;
    }

    public void createFridgeItems(FridgeItemListDTO items) {
        List<FridgeItemEntity> entities = items.fridgeItems().stream()
                .map(dto -> fridgeItemRepository.findByName(dto.name())
                        .map(existing -> {
                            existing.setQuantity(dto.quantity());
                            existing.setExpirationDate(dto.expirationDate());
                            return existing;
                        })
                        .orElse(FridgeItemMapper.toEntity(dto)))
                .toList();
        fridgeItemRepository.saveAll(entities);
    }
}
