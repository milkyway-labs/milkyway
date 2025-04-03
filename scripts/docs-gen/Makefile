venv:
    # Create a Python virtual environment if it doesn't exist
	@if [ ! -d venv ]; then python -m venv venv; fi

install-requirements: venv
	. "${CURDIR}/venv/bin/activate"; pip install -r requirements.txt 

gen: venv
	. "${CURDIR}/venv/bin/activate"; python docs-gen.py "${CURDIR}/../../x" "${CURDIR}/../../docs/generated/x"

release: venv
	. "${CURDIR}/venv/bin/activate"; python docs-release.py "${CURDIR}/../../x"
