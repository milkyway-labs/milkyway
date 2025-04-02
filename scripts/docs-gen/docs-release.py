import os
import shutil
import argparse
from packaging.version import Version
from utils import generate_docs
import re


def generate_release_docs(
        modules_dir: str,
        docs_dir: str,
        release_version: str,
):
    print(f"Modules directory: {modules_dir}")
    print(f"Docs directory: {docs_dir}")
    print(f"Preparing documentation: {release_version}")

    # Remove the old documentation
    version_dir = os.path.join(docs_dir, release_version)
    if os.path.exists(version_dir):
        shutil.rmtree(version_dir)

    # Generate the documentation
    generate_docs(modules_dir, version_dir, True)

    # Function to filter out the version dirs
    def version_dirs_filter(dir_name: str) -> bool:
        if dir_name == "main":
            return False
        dir_path = os.path.join(docs_dir, dir_name)
        return os.path.isdir(dir_path)
    # Get the various versions
    versions = list(filter(version_dirs_filter, os.listdir(docs_dir)))
    # Sort the versions in descending order
    versions.sort(key=Version, reverse=True)

    # Remove the old README.md
    modules_readme = os.path.join(docs_dir, "README.md")
    if os.path.exists(modules_readme):
        os.remove(modules_readme)


def update_summary(summary_file: str, docs_dir: str):
    summary = {}  # type: Dict[str, List[str]]
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
            summary[module] = []
            continue

        # We may have nested folders inside the module's folder
        # this allows to ignore those cases
        version_split = version.split("/")
        version = version_split[1]
        if version not in summary:
            continue

        # Add the module to the summary
        summary[version].append(module)

    # Generate the new summary
    new_content = []
    for version, modules in summary.items():
        new_content.append(f"* [{version}](modules/{version}/README.md)")
        modules.sort()
        for module in modules:
            new_content.append(f"  * [x/{module}](modules/{version}/{module}/README.md)")

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
    print(updated_content)

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
    # generate_release_docs(args.modules, docs_dir, release_version)

    # Update the Gitbook summary
    update_summary(summary_file, docs_dir)


if __name__ == "__main__":
    main()
