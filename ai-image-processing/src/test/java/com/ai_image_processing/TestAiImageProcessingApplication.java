package com.ai_image_processing;

import org.springframework.boot.SpringApplication;

public class TestAiImageProcessingApplication {

	public static void main(String[] args) {
		SpringApplication.from(AiImageProcessingApplication::main).with(TestcontainersConfiguration.class).run(args);
	}

}
