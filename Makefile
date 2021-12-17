include .env
export $(shell sed 's/=.*//' .env)

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
	pipenv run python main.py

configure-cron:
	(crontab -l ; echo "0 8 * * * cd $${PWD} && python3 main.py") | sort - | uniq - | crontab -

create-db:
	sqlite3 ${DB_LOCATION} \
		'CREATE TABLE frser_sqlite (timestamp INTEGER, version TEXT, downloads INTEGER)'

install-configure: configure-cron create-db
