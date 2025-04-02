import requests
import os
import sys
import re
import io


def github_url_to_raw(url: str) -> str:
    return url.replace("github.com", "raw.githubusercontent.com").replace("/blob/", "/")


def get_lines(text: str, start: int, end: int):
    lines = text.split('\n')  # Split the string into lines
    return '\n'.join(lines[start-1:end])


def write_gitbook_meta(file: io.TextIOWrapper):
    file.write("""---
layout:
  title:
    visible: true
  description:
    visible: false
  tableOfContents:
    visible: true
  outline:
    visible: false
  pagination:
    visible: true
---

""")


def process_markdown_file(file_path: str, output_dir: str, gitbook_meta: bool):
    print(f"Processing {file_path}")
    with open(file_path, 'r') as file:
        content = file.read()

    # Regular expression pattern for matching the section
    pattern = r'\`\`\`\w+ reference\n((?:https://)?github\.com/.*?)\n\`\`\`'
    references = re.findall(pattern, content)

    for reference_url in set(references):  # Use set to remove duplicates
        values = reference_url.split('#L')
        if len(values) != 2:
            print(f"Invalid link: {reference_url} in file ${file_path}")
            continue

        # Get the url and the lines
        url, lines = values
        start_end_line = lines.replace('L', '').split('-')
        if len(start_end_line) != 2:
            print(f"Invalid line format: {lines} in file ${file_path}")
            continue
        start, end = start_end_line

        # Parse the lines
        start = int(start)
        end = int(end)
        if start < 0 or end < 0 or start >= end:
            print(f"Invalid line range: {start}-{end} in file ${file_path}")
            continue
        # Convert the reference url to a raw github url
        raw_url = github_url_to_raw(url)

        # Extract the section from the remote reference
        res = requests.get(raw_url, allow_redirects=True)
        reference_text = get_lines(res.text, start, end)

        # Update the content with the reference
        content = content.replace(
            f" reference\n{reference_url}", f"\n{reference_text}")

    # Save the file
    parent = os.path.dirname(output_dir)
    if not os.path.exists(parent):
        os.makedirs(parent)

    with open(output_dir, 'w') as file:
        if gitbook_meta:
            write_gitbook_meta(file)
        file.write(content)
    print(f"Saved to {output_dir}")


def generate_docs(modules_dir: str, output_dir: str, gitbook_meta: bool = False):
    # Check if the modules directory exists
    if not os.path.isdir(modules_dir):
        sys.exit('Error: The modules directory does not exist.')

    # If output directory does not exist, create it
    if not os.path.exists(output_dir) or not os.path.isdir(output_dir):
        os.makedirs(output_dir)

    # Recursively search for markdown files in the modules directory
    for foldername, subfolders, filenames in os.walk(modules_dir):
        # Prevent analysis of .md files in sub directories
        relative_folder = foldername.replace(modules_dir, "")
        relative_path_split = relative_folder.split("/")
        if len(relative_path_split) > 2:
            continue

        for filename in filenames:
            if filename.endswith('.md'):
                print(f"Processing {filename} in {foldername}")
                file_path = os.path.join(foldername, filename)
                output_path = file_path.replace(modules_dir, output_dir)
                process_markdown_file(file_path, output_path, gitbook_meta)

    # Delete the README.md file if it exists
    readme_file = os.path.join(output_dir, 'README.md')
    if os.path.exists(readme_file):
        os.remove(readme_file)

    # Prepare the list of modules
    modules = os.listdir(output_dir)
    modules.sort()

    # Generate the README.md file
    with open(readme_file, 'w') as readme:
        if gitbook_meta:
            write_gitbook_meta(readme)
        readme.write("# Modules\n\n")
        for module in modules:
            readme.write(f"* [{module}]({module}/README.md)\n")
