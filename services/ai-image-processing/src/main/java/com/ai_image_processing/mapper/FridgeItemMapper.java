package com.ai_image_processing.mapper;

import com.ai_image_processing.dto.FridgeItemDTO;
import com.ai_image_processing.models.FridgeItemEntity;

public class FridgeItemMapper {

    private FridgeItemMapper(){}

    public static FridgeItemEntity toEntity(FridgeItemDTO dto){
        FridgeItemEntity entity = new FridgeItemEntity();
        entity.setName(dto.name());
        entity.setQuantity(dto.quantity());
        entity.setExpirationDate(dto.expirationDate());
        return entity;
    }
}
