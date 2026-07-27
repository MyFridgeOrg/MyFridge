import base64
import re
from pathlib import Path

import litellm

from app.schemas.fridge_item import FridgeItemListSchema

SYSTEM_PROMPT = (
    "Analyze the image of the inside of a fridge. "
    "Identify ALL visible food items, not just one. "
    "For each item, estimate the quantity as a number and expiration date if visible. "
    "If expiration date is not visible, use 'unknown'. "
    "Respond with ONLY raw JSON, no markdown, no code fences, no explanation text before or after. "
    "The JSON must match exactly this shape: "
    '{"fridge_items": [{"name": "string", "quantity": "string", "expiration_date": "string or null"}]}'
)

USER_PROMPT = "List all the food items in this fridge image."


def load_image_from_disk(filename: str) -> bytes:
    image_path = Path(__file__).resolve().parent.parent / "resources" / "images" / filename
    return image_path.read_bytes()


def _extract_json(content: str) -> str:
    """Strip markdown code fences if the model wraps its JSON response in one."""
    match = re.search(r"```(?:json)?(.*?)```", content, re.DOTALL)
    return match.group(1).strip() if match else content.strip()


def process_image(image_bytes: bytes, content_type: str) -> FridgeItemListSchema:
    base64_image = base64.b64encode(image_bytes).decode("utf-8")
    response = litellm.completion(
        model="bedrock/eu.mistral.pixtral-large-2502-v1:0",
        messages=[
            {"role": "system", "content": SYSTEM_PROMPT},
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": USER_PROMPT},
                    {"type": "image_url", "image_url": {"url": f"data:{content_type};base64,{base64_image}"}},
                ],
            },
        ],
    )
    content = response.choices[0].message.content
    if content is None:
        raise ValueError("Model returned no content")
    return FridgeItemListSchema.model_validate_json(_extract_json(content))
