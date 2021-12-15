# Grafana Plugin Loader

[![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/psf/black)

This repo is meant to house a small hobby script that downloads information about a grafana
plugin daily and stores it in a SQLite database.
This data can then be charted later in Grafana (with the SQLite plugin)

## Development

`make` is used as a task runner. For the commands take a look at the `Makefile`.

### Requirements

- make
- python 3.7 (the target raspberry runs Python 3.7.3)
- pip packages
  - pipenv

### Setup

Either the command `make install-dev` gets the dependencies or shows what is missing.

### Testing

Tests can be run with the command `make test`

### Running the script

The command `make run` should run the script and add an entry to the database.

### Release

Run `make install-production` to install dependencies.
In case this is the initial installation some further steps (like adding the script to the crontab)
should be done via `make install-configure`.
