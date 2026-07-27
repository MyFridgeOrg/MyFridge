from fastapi import FastAPI

from app.routes.fridge import router as fridge_router

app = FastAPI(title="AI Image Processing Service")

app.include_router(fridge_router)

@app.get("/health")
async def health():
    return {"status": "ok"}