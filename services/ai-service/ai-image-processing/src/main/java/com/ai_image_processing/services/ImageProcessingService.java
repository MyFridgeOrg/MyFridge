package com.ai_image_processing.services;

import com.ai_image_processing.dto.FridgeItemListDTO;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.util.MimeType;
import org.springframework.util.MimeTypeUtils;

import java.io.IOException;
import java.io.InputStream;

@Service
public class ImageProcessingService {

    private final ChatClient chatClient;

    public ImageProcessingService(ChatClient.Builder chatClientBuilder) {
        this.chatClient = chatClientBuilder.build();
    }

    public byte[] loadImageFromResources(String filename) throws IOException {
        ClassPathResource resource = new ClassPathResource(filename);
        try (InputStream inputStream = resource.getInputStream()) {
            return inputStream.readAllBytes();
        }
    }

    public FridgeItemListDTO processImage(byte[] imageBytes, String contentType) {
        try {
            Resource imageResource = new ByteArrayResource(imageBytes);
            MimeType mimeType = MimeTypeUtils.parseMimeType(contentType);

           return chatClient.prompt()
                    .system(systemMessage -> systemMessage
                            .text("Analyze the image of the inside of a fridge")
                            .text("Identify ALL visible food items, not just one")
                            .text("For each item, estimate the quantity as a number and expiration date if visible")
                            .text("If expiration date is not visible, use 'unknown'")
                    )
                    .user(userMessage -> userMessage
                            .text("List all the food items in this fridge image")
                            .media(mimeType, imageResource)
                    )
                    .call()
                    .entity(FridgeItemListDTO.class);
        } catch (Exception e) {
            throw new RuntimeException("Failed to process image", e);
        }
    }
}
