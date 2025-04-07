import argparse
from utils import generate_docs


def main():
    parser = argparse.ArgumentParser(
        description="Check module directory and output directory.")
    parser.add_argument("modules", help="The modules directory")
    parser.add_argument("output", help="The output directory")
    args = parser.parse_args()
    generate_docs(args.modules, args.output)


if __name__ == "__main__":
    main()
