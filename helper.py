import os
import re
import shutil
import sys


def to_camel_case(name: str) -> str:
    return name[0].lower() + name[1:] if name else name


def to_snake_case(name: str) -> str:
    return re.sub(r"(?<!^)(?=[A-Z])", "_", name).lower()


def replace_in_file(file_path: str, replacements: dict):
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    for old, new in replacements.items():
        content = content.replace(old, new)

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)


def main(old_name: str, new_name: str):
    old_name_camel = to_camel_case(old_name)
    new_name_camel = to_camel_case(new_name)
    old_name_snake = to_snake_case(old_name)
    new_name_snake = to_snake_case(new_name)

    replacements = {
        old_name: new_name,
        old_name_camel: new_name_camel,
        old_name_snake: new_name_snake,
    }

    for filename in os.listdir():
        if filename.endswith(".go") and old_name_snake in filename:
            new_filename = filename.replace(old_name_snake, new_name_snake)
            shutil.copy(filename, new_filename)
            replace_in_file(new_filename, replacements)
            print(f"Processed: {new_filename}")


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <OldName> <NewName>")
        sys.exit(1)

    old_name = sys.argv[1]
    new_name = sys.argv[2]
    main(old_name, new_name)
