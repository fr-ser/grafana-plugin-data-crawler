import json
import logging
from logging.handlers import RotatingFileHandler
import os
import sqlite3
import time
from urllib.request import urlopen

DB_LOCATION = os.environ.get("DB_LOCATION", "./plugin.db")
GRAFANA_PLUGIN_URL = "https://grafana.com/api/plugins/frser-sqlite-datasource/versions"


def setup_logging():
    for logger in (logging.getLogger(name) for name in logging.root.manager.loggerDict):
        logger.setLevel(logging.INFO)

    logging.basicConfig(
        handlers=[
            RotatingFileHandler("./app.log", maxBytes=5_000_000, backupCount=1),
            logging.StreamHandler(),
        ],
        level=logging.DEBUG,
        format="[%(asctime)s] %(levelname)s | %(name)s | %(message)s",
        datefmt="%Y-%m-%dT%H:%M:%S%z",
    )


def get_plugin_data():
    with urlopen(GRAFANA_PLUGIN_URL) as response:
        plugin_data = json.loads(response.read().decode("utf-8"))
    return plugin_data


def save_to_db(plugin_data):
    now_in_seconds = time.time()
    con = sqlite3.connect(DB_LOCATION)
    cur = con.cursor()

    insert_data = [
        (now_in_seconds, item["version"], item["downloads"]) for item in plugin_data["items"]
    ]

    cur.executemany(
        "INSERT INTO frser_sqlite (timestamp, version, downloads) VALUES (?, ?, ?)", insert_data
    )
    con.commit()
    con.close()


def main():
    logging.info("app started")
    plugin_data = get_plugin_data()
    logging.info(f"downloaded data with {len(plugin_data['items'])} items")
    save_to_db(plugin_data)
    logging.info("inserted the data into the database")
    logging.info("app finished")


if __name__ == "__main__":
    setup_logging()
    try:
        main()
    except Exception:
        logging.exception("top level exception")
        raise
