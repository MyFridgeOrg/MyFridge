from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from app.models.fridge_item import FridgeItem
from app.schemas.fridge_item import FridgeItemListSchema


async def upsert_fridge_items(session: AsyncSession, items: FridgeItemListSchema):
    for item in items.fridge_items:
        result = await session.execute(
            select(FridgeItem).where(FridgeItem.name.ilike(item.name)) # Match name case-insensitive
        )
        existing_item = result.scalar_one_or_none() #NOTE: We expect one or zero result. Multiple results indicate a bug and raise an error.

        if existing_item:
            existing_item.quantity = item.quantity
            existing_item.expiration_date = item.expiration_date
        else:
            session.add(
                FridgeItem(
                    name=item.name,
                    quantity=item.quantity,
                    expiration_date=item.expiration_date
                )
            )
    await session.commit()