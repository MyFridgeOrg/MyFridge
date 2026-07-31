import ast
import csv

def load_recipes(path, limit):
    with open(path, newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        for i, row in enumerate(reader):
            if i >= limit:
                break
            title = row["name"]
            ingredients = ast.literal_eval(row["ingredients"])
            steps = ast.literal_eval(row["steps"])
            instructions = ast.literal_eval(row["instructions"])
            yield title, ingredients, steps, instructions