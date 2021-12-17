# Grafana Plugin Loader

[![CI](https://github.com/fr-ser/grafana-plugin-data-crawler/actions/workflows/ci_cd.yml/badge.svg)](https://github.com/fr-ser/grafana-plugin-data-crawler/actions/workflows/ci_cd.yml)

This repo is meant to house a small hobby script that downloads information about a Grafana
plugin daily and stores it in a SQLite database.
This data can then be charted later in Grafana (with the SQLite plugin)

## Development

`make` is used as a task runner. For the commands take a look at the `Makefile`.

### Requirements

- make
- go 1.17

### Setup

Either the command `make install-dev` gets the dependencies or shows what is missing.

### Testing

Tests can be run with the command `make test`

### Running the script

The command `make run` should run the script and add an entry to the database.

### Release

1. Run `make build`
2. Copy the executable over to the target folder/device.
3. In case this is the initial installation some further steps (like adding the script to the
   crontab) should be done via make install-configure.
