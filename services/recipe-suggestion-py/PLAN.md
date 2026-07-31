# Recipe Suggestion (RAG) — Plan

## Context

`ai-image-processing-py` detects what's currently in the fridge and is now the
canonical system-of-record for fridge contents (the Java `ai-image-processing`
service and its Pi → S3 → RabbitMQ pipeline are being retired; that pipeline
will eventually be rebuilt around the Python service instead).

This new service, `recipe-suggestion-py`, answers the question **"what can I
cook tonight with what's in my fridge?"** — reframed as a chef surveying
available ingredients — using Retrieval-Augmented Generation (RAG) over a
recipe corpus, grounded and prioritized by what's about to expire.

---

## Why RAG (and not just prompting an LLM)

A pure "stuff the fridge contents into a prompt and ask for a recipe" approach
works, but only leans on whatever recipes the LLM already memorized during
training. RAG earns its keep here because:

- It grounds suggestions in a **real, curated recipe corpus** (verified steps,
  ingredients, not hallucinated), and is extensible later to the user's own
  saved/private recipes — content an LLM could never have been trained on.
- Ingredient matching against "what's actually in the fridge" is largely a
  **structured filtering problem**, not a semantic one — vector search is
  reserved for the part it's actually good at: fuzzy matching (synonyms,
  typos, mood-based queries like "something quick and vegetarian").

---

## Architecture Decisions

| Concern | Decision |
|---|---|
| Use case | Recipe suggestions from current fridge contents, prioritized by near-expiry items |
| Corpus | Public recipe dataset (e.g. RecipeNLG / Food.com-Kaggle), capped to ~10,000 recipes initially |
| Vector store | `pgvector` extension on the existing Postgres instance (no new infra) |
| Embedding model | Bedrock Titan Embeddings, called via `litellm` (consistent with existing stack) |
| Retrieval strategy | Staged: hard ingredient-overlap % filter → vector-similarity rank among survivors |
| Fallback (no filter matches) | Fall back to pure vector search, and explicitly surface which ingredients are missing / that the overlap threshold wasn't met |
| Ingestion | Manual one-off CLI script (e.g. `python -m app.ingest_recipes`), run on-demand — not an Alembic migration, not auto-run on startup |
| Schema management | Own Alembic setup for this service's tables (`recipes`, embeddings column), separate from `ai-image-processing-py`'s migrations |
| API style | Synchronous REST (matches existing service conventions; latency is a few seconds) |
| Fridge data source | `ai-image-processing-py` is canonical for current fridge contents |
| Data access pattern | Calls `ai-image-processing-py`'s `GET /fridge/items` over REST — **not** a direct DB read — to preserve a real service boundary |
| Location | `services/recipe-suggestion-py` (flat, alongside other services) |

---

## Why staged filter-then-rank instead of a weighted blend

Ingredient-overlap count and vector cosine-distance live on incompatible
numeric scales. A weighted blend (`overlap * 0.7 + similarity * 0.3`) requires
normalizing both before the weights are meaningful, and tuning those weights
by feel is a rabbit hole. Staged filtering avoids this entirely:

1. **Hard filter**: keep only recipes where ingredient overlap ≥ a threshold
   (e.g. "at least 50% of ingredients present").
2. **Rank survivors** purely by vector similarity — letting vector search do
   only the job it's good at (ranking/tie-breaking among already-plausible
   candidates), not blended arithmetic with an incompatible score.

If nothing survives the filter, fall back to pure vector search over the full
corpus and explicitly tell the user what's missing, rather than returning
nothing.

---

## High-Level Flow

**Ingestion (offline, manual, one-off):**
1. Load the capped recipe dataset (title, ingredients, instructions).
2. Embed each recipe (title + ingredients) via Bedrock Titan Embeddings.
3. Store `{recipe_id, title, ingredients, instructions, embedding}` in the
   `recipes` table (pgvector column) via the CLI ingestion script.

**Query (synchronous, per request):**
1. Call `ai-image-processing-py`'s `GET /fridge/items` to get current fridge
   contents (`name`, `quantity`, `expiration_date`).
2. Hard-filter the recipe corpus by ingredient overlap %.
3. Embed a query representation of the fridge context and rank survivors by
   vector similarity (fall back to full-corpus vector search if the filter
   yields nothing, noting missing ingredients).
4. Augment an LLM prompt with the top-K retrieved recipes + fridge context
   (prioritizing near-expiry items).
5. Generate and return a tailored recommendation via a new endpoint, e.g.
   `GET /recipes/suggestions`.

---

## Open Items / Not Yet Decided

- Exact overlap-% threshold for the hard filter (needs empirical tuning).
- Recipe dataset source finalization (RecipeNLG vs. Food.com/Kaggle vs. other).
- Whether/when to fold in personal/saved recipes as a second corpus source.
- Fate of the Java `ai-image-processing` service and its CI workflows, now
  that Python is canonical (separate decision, not addressed here).
