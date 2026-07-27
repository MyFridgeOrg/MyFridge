package com.ai_image_processing.event;

import java.io.IOException;
import java.util.function.Consumer;

import com.ai_image_processing.dto.FridgeItemListDTO;
import com.ai_image_processing.services.FridgeItemService;
import com.ai_image_processing.services.ImageProcessingService;
import com.ai_image_processing.services.S3ImageService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ImageEventConsumer {

    private static final Logger log = LoggerFactory.getLogger(ImageEventConsumer.class);

    @Value("${spring.application.test-image}")
    private String testImagePath;

    @Bean
    public Consumer<ImageCapturedMessage> handleImageCaptured(
            S3ImageService s3ImageService,
            ImageProcessingService imageProcessingService,
            FridgeItemService fridgeItemService
    ) {
        return message -> {
            log.info("Received image event: device={}, object={}",
                    message.deviceId(), message.objectKey());
            byte[] imageBytes = null;
            try {
                imageBytes = imageProcessingService.loadImageFromResources(testImagePath);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
//          byte[] imageBytes = s3ImageService.download(message.bucket(), message.objectKey());
            FridgeItemListDTO items = imageProcessingService.processImage(imageBytes, message.contentType());
            fridgeItemService.createFridgeItems(items);
            log.info("Successfully processed image: {}", message.objectKey());
        };
    }
}

