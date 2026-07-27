# RabbitMQ Consumer Implementation Plan

## Context

The Go application running on a Raspberry Pi captures fridge images, uploads them to S3, and publishes a message to RabbitMQ. The Spring Boot `ai-image-processing` service consumes that message, fetches the image from S3, processes it with an AI model, and persists the extracted fridge items to PostgreSQL.

---

## Architecture Decisions

| Concern | Decision |
|---------|----------|
| Broker | RabbitMQ (cloud-agnostic, runs anywhere) |
| Storage | S3 (Go app abstracts behind an interface for portability) |
| Spring integration | Spring Cloud Stream + RabbitMQ binder |
| S3 client | Spring Cloud AWS S3 starter (`S3Template`) |
| Topology | Exchange: `myfridge.events` (topic), routing key: `image.captured`, queue: `myfridge.events.ai-image-processing` |
| Message payload | `{ bucket, objectKey, deviceId, capturedAt, contentType }` |
| Consumer role | Pure consumer — no publishing |
| Error handling | 3 retries with exponential backoff → dead-letter queue |
| Local dev | Arconia Dev Services for RabbitMQ (auto-starts container) |
| REST endpoint | Keep `POST /run` for dev/testing |

---

## Project Structure

```
src/main/java/com/ai_image_processing/
├── controller/
│   └── ImageProcessingController.java        # Thin REST adapter (dev/test trigger)
├── event/
│   ├── ImageCapturedMessage.java             # Record DTO for the incoming message
│   └── ImageEventConsumer.java              # @Configuration with Consumer<> bean (orchestrates the full flow)
├── services/
│   ├── ImageProcessingService.java           # Calls LLM with image bytes, returns FridgeItemListDTO
│   ├── FridgeItemService.java                # Persists fridge items to DB (existing)
│   └── S3ImageService.java                   # Fetches image bytes from S3 via S3Template
├── dto/
│   ├── FridgeItemDTO.java                    # Existing
│   └── FridgeItemListDTO.java                # Existing
├── mapper/
│   └── FridgeItemMapper.java                 # Existing
├── models/
│   └── FridgeItemEntity.java                 # Existing
└── repository/
    └── FridgeItemRepository.java             # Existing
```

**Note:** No separate orchestration service. The consumer bean in `ImageEventConsumer` directly coordinates the flow:
1. Consume message from queue
2. Download image from S3
3. Process image with LLM
4. Save results to database

This avoids an unnecessary abstraction layer — the glue logic is trivial (3 service calls in sequence).

---

## Message Contract

```json
{
  "bucket": "my-fridge-images",
  "objectKey": "captures/2026-06-19/pi-001/image-1718793600.jpg",
  "deviceId": "pi-001",
  "capturedAt": "2026-06-19T12:00:00Z",
  "contentType": "image/jpeg"
}
```

---

## RabbitMQ Topology

- **Exchange:** `myfridge.events` (type: topic, durable)
- **Routing key:** `image.captured`
- **Queue:** `myfridge.events.ai-image-processing` (durable, bound to exchange with routing key)
- **DLQ:** `myfridge.events.ai-image-processing.dlq`

The Go app publishes to the exchange with routing key `image.captured`. Spring Cloud Stream creates the queue and binds it automatically.

---

## Implementation Steps

### 1. Add Gradle Dependencies

```groovy
// Spring Cloud Stream + RabbitMQ binder
implementation 'org.springframework.cloud:spring-cloud-stream-binder-rabbit'

// Spring Cloud AWS S3
implementation 'io.awspring.cloud:spring-cloud-aws-starter-s3'

// Arconia Dev Services for RabbitMQ (local dev)
testAndDevelopmentOnly 'io.arconia:arconia-dev-services-rabbitmq'

// Spring Cloud Stream test binder (for integration tests)
testImplementation("org.springframework.cloud:spring-cloud-stream") {
    artifact {
        name = "spring-cloud-stream"
        extension = "jar"
        type = "test-jar"
        classifier = "test-binder"
    }
}
```

Add Spring Cloud BOM for version management.

### 2. Create `ImageCapturedMessage`

```java
// event/ImageCapturedMessage.java
public record ImageCapturedMessage(
    String bucket,
    String objectKey,
    String deviceId,
    String capturedAt,
    String contentType
) {}
```

### 3. Create `S3ImageService`

```java
// services/S3ImageService.java
@Service
public class S3ImageService {

    private final S3Template s3Template;

    public S3ImageService(S3Template s3Template) {
        this.s3Template = s3Template;
    }

    public byte[] download(String bucket, String objectKey) {
        Resource resource = s3Template.download(bucket, objectKey);
        try {
            return resource.getContentAsByteArray();
        } catch (IOException e) {
            throw new RuntimeException("Failed to download image from S3", e);
        }
    }
}
```

### 4. Refactor `ImageProcessingService`

Change `processImage` to accept `byte[]` and a MIME type string instead of a classpath filename:

```java
// services/ImageProcessingService.java
public FridgeItemListDTO processImage(byte[] imageBytes, String contentType) {
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
}
```

### 5. Create `ImageEventConsumer`

```java
// event/ImageEventConsumer.java
@Configuration
public class ImageEventConsumer {

    private static final Logger log = LoggerFactory.getLogger(ImageEventConsumer.class);

    @Bean
    public Consumer<ImageCapturedMessage> handleImageCaptured(
            S3ImageService s3ImageService,
            ImageProcessingService imageProcessingService,
            FridgeItemService fridgeItemService
    ) {
        return message -> {
            log.info("Received image captured event: device={}, object={}",
                    message.deviceId(), message.objectKey());
            byte[] imageBytes = s3ImageService.download(message.bucket(), message.objectKey());
            FridgeItemListDTO items = imageProcessingService.processImage(imageBytes, message.contentType());
            fridgeItemService.createFridgeItems(items);
            log.info("Successfully processed image: {}", message.objectKey());
        };
    }
}
```

### 6. Update `ImageProcessingController`

Keep for dev/testing — loads a classpath image and calls services directly:

```java
@PostMapping("/run")
public ResponseEntity<?> run() {
    byte[] imageBytes = loadImageFromResources(testImagePath);
    FridgeItemListDTO items = imageProcessingService.processImage(imageBytes, "image/png");
    fridgeItemService.createFridgeItems(items);
    return ResponseEntity.ok().build();
}
```

### 7. Configure `application.yml`

```yaml
spring:
  cloud:
    function:
      definition: handleImageCaptured
    stream:
      bindings:
        handleImageCaptured-in-0:
          destination: myfridge.events
          group: ${spring.application.name}
      rabbit:
        bindings:
          handleImageCaptured-in-0:
            consumer:
              bindingRoutingKey: image.captured
              exchangeType: topic
              maxAttempts: 3
              backOffInitialInterval: 1000
              backOffMultiplier: 2.0
              autoBindDlq: true
              dlqTtl: 86400000
  rabbitmq:
    host: localhost
    port: 5672
    username: guest
    password: guest
    connection-timeout: 5s
  cloud.aws:
    s3:
      region: eu-central-1
```

### 8. Arconia Dev Services Config

Arconia should auto-detect the RabbitMQ dependency and start a container. If explicit config is needed:

```yaml
arconia:
  dev:
    services:
      rabbitmq:
        shared: true
```

---

## Networking (TBD)

The Raspberry Pi needs to reach the RabbitMQ broker. Options to evaluate:

- **Tailscale** mesh VPN (recommended — private network, no public exposure)
- **Public endpoint** with TLS + credentials
- **WireGuard** tunnel

This needs resolution before deployment.

---

## Go Application Side (for reference)

The Go app will:
1. Capture image
2. Upload to S3 (behind a storage interface)
3. Publish message to RabbitMQ exchange `myfridge.events` with routing key `image.captured`

Message published by Go:
```json
{
  "bucket": "my-fridge-images",
  "objectKey": "captures/2026-06-19/pi-001/image-1718793600.jpg",
  "deviceId": "pi-001",
  "capturedAt": "2026-06-19T12:00:00Z",
  "contentType": "image/jpeg"
}
```

