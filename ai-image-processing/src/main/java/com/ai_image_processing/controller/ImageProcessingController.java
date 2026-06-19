package com.ai_image_processing.controller;

import com.ai_image_processing.dto.FridgeItemListDTO;
import com.ai_image_processing.services.FridgeItemService;
import com.ai_image_processing.services.ImageProcessingService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ImageProcessingController {

    @Value("${spring.application.test-image}")
    private String testImagePath;

    private final ImageProcessingService imageProcessingService;
    private final FridgeItemService fridgeItemService;

    public ImageProcessingController(ImageProcessingService imageProcessingService, FridgeItemService fridgeItemService) {
        this.imageProcessingService = imageProcessingService;
        this.fridgeItemService = fridgeItemService;
    }

    @PostMapping("/run")
    public ResponseEntity<?> run(
    ){
        FridgeItemListDTO items = imageProcessingService.processImage(testImagePath);
        fridgeItemService.createFridgeItems(items);
        return ResponseEntity.ok().build();
    }
}
