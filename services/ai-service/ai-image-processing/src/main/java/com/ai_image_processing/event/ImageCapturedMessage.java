package com.ai_image_processing.event;

public record ImageCapturedMessage(
        String bucket,
        String objectKey,
        String deviceId,
        String capturedAt,
        String contentType
) {}
