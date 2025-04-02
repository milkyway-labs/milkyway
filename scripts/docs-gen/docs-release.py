import os
import shutil
import argparse
from packaging.version import Version
from utils import generate_docs, write_gitbook_meta


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

    # Generate the new README.md
    with open(modules_readme, "w") as readme:
        write_gitbook_meta(readme)
        readme.write("# Modules version\n\n")

        # Print first the main version
        if os.path.exists(os.path.join(docs_dir, "main")):
            readme.write("- [main](main/README.md)\n")
        # Print the various versions
        for version in versions:
            readme.write(f"- [{version}]({version}/README.md)\n")


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

    generate_release_docs(args.modules, docs_dir, release_version)


if __name__ == "__main__":
    main()
