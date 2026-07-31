import pgvector.sqlalchemy
from sqlalchemy import ARRAY, String, Text
from sqlalchemy.orm import Mapped, mapped_column
from pgvector.sqlalchemy import Vector

from app.config.database import Base

class Recipe(Base):
    __tablename__ = 'recipes'

    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(255), nullable=False)
    ingredients: Mapped[list[str]] = mapped_column(ARRAY(String), nullable=False)
    instructions: Mapped[str] = mapped_column(Text, nullable=False)
    embedding: Mapped[list[float]] = mapped_column(Vector(1024), nullable=False)