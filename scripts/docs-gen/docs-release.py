import os
import argparse
from utils import generate_docs
import re
import pathlib
from typing import List


def get_modules(docs_dir: str) -> List[str]:
    """
    Gets the list of modules for which exists the documentation 
    inside the provided directory
    """
    modules = []  # type: List[str]
    # Generate the summary section
    for foldername, _, _ in os.walk(docs_dir):
        folder = foldername.replace(docs_dir, "")
        # Ignore the docs dir
        if folder == "":
            continue
        # Ignore path deeper then 2 and the version path
        folder_split = folder.split(os.path.sep)
        if len(folder_split) < 2:
            continue
        _, module = folder_split

        # Add the module to the summary
        modules.append(module)

    modules.sort()
    return modules

def generate_modules_list(modules: List[str], parent_dir: str = '') -> List[str]:
    """
    Generates a markdown list of the provided modules.
    """
    modules_list = []
    if parent_dir != '' and parent_dir[-1] != '/':
        parent_dir += '/'

    for module in modules:
        modules_list.append(f"* [x/{module}]({parent_dir}{module}/README.md)")

    return modules_list 



def generate_release_docs(
        modules_dir: str,
        docs_dir: str,
):
    """
    Generates the documentation for the provided release version.
    """
    print(f"Modules directory: {modules_dir}")
    print(f"Docs directory: {docs_dir}")

    # Generate the documentation
    generate_docs(modules_dir, docs_dir, True)


def generate_modules_readme(docs_dir: str):
    """
    Generates the README.md inside the directory that contains the modules  documentation.
    """
    # Remove the old README.md
    modules_readme = os.path.join(docs_dir, "README.md")
    if os.path.exists(modules_readme):
        os.remove(modules_readme)

    # Gets the modules inside the documentation dir
    modules = get_modules(docs_dir)

    # Generate the final modules/README.md file using
    # the modules-template.md file as a template
    script_path = pathlib.Path(__file__).parent.resolve()
    template_file_file = os.path.join(script_path, 'modules-template.md')
    modules_template = open(template_file_file, 'r').read()
    modules_file = modules_template.replace('{{ modules }}', "\n".join(generate_modules_list(modules)))

    # Writes the file
    with open(modules_readme, 'w', encoding='utf-8') as file:
        file.write(modules_file)


def update_summary(summary_file: str, docs_dir: str):
    """
    Updates the Gitbook SUMMARY.md file to include the generated documentation
    present inside the modules directory.
    """
    modules = get_modules(docs_dir)
    # Generate the new summary
    new_content = generate_modules_list(modules, 'modules')

    # Update the summary file
    with open(summary_file, 'r', encoding='utf-8') as file:
        content = file.read()

    # Regular expression to match content between the <!-- modules --> tags,
    # ensuring it works even if the section is empty or has indentation
    # Matches everything between the tags
    pattern = r'(\s*)(<!-- modules -->\n)(.*?)(\s*<!-- modules -->)'

    # Matches everything between the tags
    def replace_match(match):
        indent, start_tag, _, end_tag = match.groups()
        indent = indent.replace('\n', '')
        end_tag = end_tag.replace('\n', '')
        indented_content = '\n'.join(
            [indent + line if line.strip() else indent for line in new_content])
        return f'\n{indent}{start_tag}{indented_content}\n{end_tag}'

    updated_content = re.sub(pattern, replace_match, content, flags=re.DOTALL)

    with open(summary_file, 'w', encoding='utf-8') as file:
        file.write(updated_content)


def main():
    parser = argparse.ArgumentParser(
        description="Generate the documentation to be published")
    parser.add_argument("modules", help="The modules directory")
    args = parser.parse_args()

    docs_dir = os.getenv("DOCS_DIR")
    if docs_dir is None:
        docs_dir = "./test"

    summary_file = os.getenv("GITBOOK_SUMMARY")
    if summary_file is None:
        summary_file = "./test/SUMMARY.md"

    # Generate the release documentation
    generate_release_docs(args.modules, docs_dir)

    # Generate a README.md inside the docs dir
    generate_modules_readme(docs_dir)

    # Update the Gitbook summary
    update_summary(summary_file, docs_dir)


if __name__ == "__main__":
    main()
