from pydantic import BaseModel

class FridgeItemSchema(BaseModel):
    name: str
    quantity: str
    expiration_date: str | None = None

class FridgeItemListSchema(BaseModel):
    fridge_items: list[FridgeItemSchema]