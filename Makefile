include .env
export $(shell sed 's/=.*//' .env)

install-dev:
	go install
	go get -u golang.org/x/lint/golint

test-lint:
	go vet .
	golint .

test-unit:
	go test

test: test-lint test-unit

run:
	go run main.go

configure-cron:
	(crontab -l ; echo "0 3 * * * cd $${PWD} && ./build/grafana-plugin-data-crawler") | sort - | uniq - | crontab -

create-db:
	sqlite3 ${DB_LOCATION} \
		'CREATE TABLE frser_sqlite (timestamp INTEGER, version TEXT, downloads INTEGER)'

install-configure: configure-cron create-db

build-local:
	@rm -f build/grafana-plugin-data-crawler
	go build -o build/grafana-plugin-data-crawler

build-raspberrypi-zero:
	@rm -f build/grafana-plugin-data-crawler_linux_arm5
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/grafana-plugin-data-crawler_linux_arm5

.PHONY: build
build: build-local build-raspberrypi-zero
