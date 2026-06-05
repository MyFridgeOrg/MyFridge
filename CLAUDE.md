# MyFridge

Smart Fridge Inventory Assistant. A Raspberry Pi detects when the fridge opens and captures an image. A Spring Boot backend receives the image, uses an AI vision model to identify food items, and exposes the results via a REST API. A React frontend displays the current fridge contents.

## Architecture

```
Raspberry Pi (Go)
  └── Detects fridge opening via sensor
  └── Captures image with camera
  └── Uploads image to backend API

Backend (Java / Spring Boot)
  └── Receives image uploads
  └── Calls AI vision model to identify food items
  └── Stores latest image and detected items
  └── Exposes REST API

Frontend (React / Vite) — future
  └── Displays latest fridge image
  └── Shows detected food items and last update time
```

## Repository Structure

```
pi/          # Go — Raspberry Pi code (camera, sensor, image upload)
backend/     # Java — Spring Boot REST API and business logic
frontend/    # React — web UI (future)
```

## Running the Backend

```bash
cd backend
./gradlew bootRun
```

## API

| Method | Path           | Description               |
|--------|----------------|---------------------------|
| GET    | /health        | Health check              |
| POST   | /images        | Upload a fridge image     |
| GET    | /images/latest | Get latest image + items  |

## Collaboration Style

This is a learning project. The developer wants to do things themselves and learn by doing — Claude's role is to explain and guide, not to execute. Concretely:

- Provide commands for the developer to run, don't run them via Bash
- Explain what a command or concept does before suggesting it
- Scaffold and write code only when asked to

## Development Phases

1. Camera proof of concept (Raspberry Pi)
2. Fridge image experiment
3. AI detection experiment
4. **Backend prototype ← current**
5. Raspberry Pi integration (Go)
6. AI integration
7. Web UI (React)
