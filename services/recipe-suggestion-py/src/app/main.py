from fastapi import FastAPI

app = FastAPI(title="Recipe Suggestion Service")


@app.get("/health")
async def health():
    return {"status": "ok"}
