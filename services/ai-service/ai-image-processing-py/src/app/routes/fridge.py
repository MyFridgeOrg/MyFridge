from typing import Annotated

from fastapi import APIRouter, Depends, UploadFile, status
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.config.database import async_session
from app.models.fridge_item import FridgeItem
from app.schemas.fridge_item import FridgeItemSchema
from app.services import fridge_item_service
from app.services.image_processing_service import process_image, load_image_from_disk

router = APIRouter()

async def get_db():
    async with async_session() as session:
        yield session

SessionDep = Annotated[AsyncSession, Depends(get_db)]

@router.get("/fridge/items")
async def get_fridge_items(session: SessionDep) -> list[FridgeItemSchema]:
    result = await session.execute(select(FridgeItem))
    items = result.scalars().all()
    return [
        FridgeItemSchema(
            name=item.name,
            quantity=item.quantity,
            expiration_date=item.expiration_date,
        )
        for item in items
    ]

@router.post("/run", status_code=status.HTTP_204_NO_CONTENT)
async def create_fridge_items(session: SessionDep) -> None:
    image_bytes = load_image_from_disk("test-image.png")
    items = process_image(image_bytes, "image/png")
    await fridge_item_service.upsert_fridge_items(session, items)
    return None
