install-dev:
	pipenv install --dev

test-unit:
	pipenv run pytest test_main.py

test-linting:
	@echo
	@pipenv run flake8 && echo "===> Linting passed! <==="
	@echo

test: test-linting test-unit

run:
	echo "not implemented" && exit 1

install-production:
	echo "not implemented" && exit 1

install-configure:
	echo "not implemented" && exit 1
