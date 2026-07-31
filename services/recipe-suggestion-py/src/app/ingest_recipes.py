import argparse

from app.services.recipe_ingestion_service import load_recipes

def parse_args():
    parser = argparse.ArgumentParser(description='Ingest recipes into vector store')
    parser.add_argument('--input', required=True, help='path to recipes file')
    parser.add_argument('--limit', type=int, default=10000, help="Max number of recipes to ingest")
    return parser.parse_args()

if __name__ == "__main__":
    args = parse_args()
    for title, ingredients, instructions in load_recipes(args.input, limit=args.limit):
        print(title)
        print(ingredients)
        print(instructions[:100], "...")
        print("---")
