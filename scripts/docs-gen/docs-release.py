import os
import shutil
import argparse
from packaging.version import Version
from utils import generate_docs
import re
import pathlib
from typing import List, Tuple


def get_modules_by_version(docs_dir: str) -> List[Tuple[str, List[str]]]:
    """
    Gets a list of the modules that each version has.
    The resulting list contains tuples where the first items is the version
    name and the second is the list of modules in that version.
    """
    modules = []  # type: List[Tuple[str, List[str]]]
    # Generate the summary section
    for foldername, subfolders, filenames in os.walk(docs_dir):
        folder = foldername.replace(docs_dir, "")
        # Ignore the docs dir
        if folder == "":
            continue
        # Ignore path deeper then 2 and the version path
        folder_split = os.path.split(folder)
        # The files are organized in version/module/README.md
        version, module = folder_split

        # We are starting exploring a new version
        if version == "/":
            modules.append([module, []])
            continue

        # Add the module to the summary
        modules[-1][1].append(module)

    # Sort the modules
    def sort_function(row: (str, List[str])) -> Version:
        if row[0] == "main":
            return Version("0")
        return Version(row[0])

    # Sort by version
    modules.sort(key=sort_function)

    # Sort the modules of each version by name
    for row in modules:
        row[1].sort()

    return modules


def generate_release_docs(
        modules_dir: str,
        docs_dir: str,
        release_version: str,
):
    """
    Generates the documentation for the provided release version.
    """
    print(f"Modules directory: {modules_dir}")
    print(f"Docs directory: {docs_dir}")
    print(f"Preparing documentation: {release_version}")

    # Remove the old documentation
    version_dir = os.path.join(docs_dir, release_version)
    if os.path.exists(version_dir):
        shutil.rmtree(version_dir)

    # Generate the documentation
    generate_docs(modules_dir, version_dir, True)


def generate_modules_readme(docs_dir: str):
    """
    Generates the README.md inside the "modules" directory that
    contains the modules documentation for each version.
    """
    # Remove the old README.md
    modules_readme = os.path.join(docs_dir, "README.md")
    if os.path.exists(modules_readme):
        os.remove(modules_readme)

    # Gets the versions inside the documentation dir
    modules_by_version = get_modules_by_version(docs_dir)

    # Generates content that will replace the {{ modules }}
    # keyword contained in the modules-template.md file.
    modules_section = []
    for (version, modules) in modules_by_version:
        modules_section.append(f"* [{version}]({version}/README.md)")
        for module in modules:
            modules_section.append(
                f"  * [x/{module}]({version}/{module}/README.md)")

    # Generate the final modules/README.md file using
    # the modules-template.md file as a template
    script_path = pathlib.Path(__file__).parent.resolve()
    template_file_file = os.path.join(script_path, 'modules-template.md')
    modules_template = open(template_file_file, 'r').read()
    modules_file = modules_template.replace(
        '{{ modules }}', "\n".join(modules_section))
    with open(modules_readme, 'w', encoding='utf-8') as file:
        file.write(modules_file)


def update_summary(summary_file: str, docs_dir: str):
    """
    Updates the Gitbook SUMMARY.md file to include the generated documentation
    present inside the modules directory.
    """
    modules_by_version = get_modules_by_version(docs_dir)
    # Generate the new summary
    new_content = []
    for (version, modules) in modules_by_version:
        new_content.append(f"* [{version}](modules/{version}/README.md)")
        modules.sort()
        for module in modules:
            new_content.append(
                f"  * [x/{module}](modules/{version}/{module}/README.md)")

    # Update the summary file
    with open(summary_file, 'r', encoding='utf-8') as file:
        content = file.read()

    # Regular expression to match content between the <!-- modules --> tags,
    # ensuring it works even if the section is empty or has indentation
    # Matches everything between the tags
    pattern = r'(\s*)(<!-- modules -->\n)(.*?)(\s*<!-- modules -->)'

    # Matches everything between the tags
    def replace_match(match):
        indent, start_tag, boh, end_tag = match.groups()
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

    release_version = os.getenv("RELEASE_VERSION")
    if release_version is None:
        release_version = "main"

    summary_file = os.getenv("GITBOOK_SUMMARY")
    if summary_file is None:
        summary_file = "./test/SUMMARY.md"

    # Generate the release documentation
    generate_release_docs(args.modules, docs_dir, release_version)

    # Generate a README.md inside the docs dir
    generate_modules_readme(docs_dir)

    # Update the Gitbook summary
    update_summary(summary_file, docs_dir)


if __name__ == "__main__":
    main()
