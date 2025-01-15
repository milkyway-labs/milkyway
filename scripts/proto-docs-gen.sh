#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
# Directory where we generate the modules documentation
MODULES_DIR="${SCRIPT_DIR}/../docs/modules"

generate_docs ()
{

    # Delete the old modules documentation
    rm -rf "${MODULES_DIR}"
    mkdir -p "${MODULES_DIR}"

    # Generate the documentation from the proto files
    buf generate --template proto/buf.gen.docs.yaml

    # Move the modules directories since buf
    # will generate the modules documentation inside a milkyway directory
    cp -r "${MODULES_DIR}/milkyway"/* "${MODULES_DIR}/"
    rm -rf "${MODULES_DIR}/milkyway"
}

# Prepare the modules index.md
echo "# Modules" > "${MODULES_DIR}/index.md"
echo "" >> "${MODULES_DIR}/index.md"


# Generate the various index files
for module_dir in $(find "${MODULES_DIR}" -mindepth 1 -maxdepth 1 -type d) ; do
    module_name=$(basename "${module_dir}")

    # Add the module information to the index file that is present inside the modules dir
    echo "- [${module_name}](./${module_name}/index.md)" >> "${MODULES_DIR}/index.md"
    
    echo "# x/${module_name} module" > "${module_dir}/index.md"
    for version_dir in $(find "${module_dir}" -mindepth 1 -maxdepth 1 -type d) ; do 
        version_name=$(basename "${version_dir}")

        echo "   - [${version_name}](./${module_name}/${version_name}/index.md)" >> "${MODULES_DIR}/index.md"
        echo "- [${version_name}](./${version_name}/index.md)" >> "${module_dir}/index.md"
    done

done
