include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build

help:
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
		| grep -v -- -- \
		| sed 'N;s/\n/###/' \
		| sed -n 's/^#: \(.*\)###\(.*\):.*/\2:###\1/p' \
		| column -t  -s '###'

#: Install go dependencies
install-dependencies:
	go mod download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

#: Run the linting tests
test-lint:
	@echo "Linting Checks:"
	@golangci-lint run ./... && echo "Linting passed!\n"

#: Run the unit tests
test-unit:
	gotestsum --format testname -- -count=1 -cover ./...

#: Run all tests
test: test-lint test-unit

#: Run the command to load data from Grafana
run-load-grafana:
	go run ./pkg/commands/load_grafana/load_grafana.go

#: Run the command to load data from Github
run-load-github:
	go run ./pkg/commands/load_github/load_github.go

#: Add the commands to the crontab
configure-cron:
	(crontab -l ; echo "0 3 * * * cd $${PWD} && ./build/load_grafana") | sort - | uniq - | crontab -
	(crontab -l ; echo "0 3 * * * cd $${PWD} && ./build/load_github") | sort - | uniq - | crontab -

#: Create the database with the required tables
create-db:
	sqlite3 ${DB_LOCATION} \
		'CREATE TABLE grafana (timestamp INTEGER, version TEXT, downloads INTEGER)'
	sqlite3 ${DB_LOCATION} \
		'CREATE TABLE github_traffic_views (timestamp INTEGER PRIMARY KEY, count INTEGER, uniques INTEGER)'

#: perform all initial steps to setup the tool on a new machine
install-configure: configure-cron create-db

#: builds executable files for the commands for the local architecture
build-local:
	@rm -f build/load_grafana
	go build -o build/load_grafana ./pkg/commands/load_grafana
	@rm -f build/load_github
	go build -o build/load_github ./pkg/commands/load_github

#: builds executable files for the commands for the raspberry pi zero architecture
build-raspberrypi-zero:
	@rm -f build/load_grafana_linux_arm5
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/load_grafana_linux_arm5 ./pkg/commands/load_grafana
	@rm -f build/load_github_linux_arm5
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/load_github_linux_arm5 ./pkg/commands/load_github

#: builds executable files for the local and respberry pi zero architecture
build: build-local build-raspberrypi-zero
