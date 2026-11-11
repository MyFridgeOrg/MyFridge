package com.ai_image_processing.services;

import io.awspring.cloud.s3.S3Template;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;

import java.io.IOException;

@Service
public class S3ImageService {

    private final S3Template s3Template;

    public S3ImageService(S3Template s3Template){
        this.s3Template = s3Template;
    }

    public byte[] download(String bucket, String objectKey){
        Resource resource = s3Template.download(bucket, objectKey);
        try {
            return resource.getContentAsByteArray();
        } catch (IOException e) {
            throw new RuntimeException("Failed to download image from S3", e);
        }
    }
}
