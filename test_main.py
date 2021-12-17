import json
import sqlite3
import time
from unittest import mock

import pytest

from main import GRAFANA_PLUGIN_URL, get_plugin_data, save_to_db


@mock.patch("main.urlopen")
def test_get_plugin_data(urlopen_mock):
    test_result = {"test": "json"}
    urlopen_mock.return_value.__enter__.return_value.read.return_value = json.dumps(
        test_result
    ).encode()

    received = get_plugin_data()

    assert urlopen_mock.called_with(GRAFANA_PLUGIN_URL)
    assert received == test_result


def test_save_to_db(tmp_path):
    test_db_path = tmp_path / "test.db"
    con = sqlite3.connect(test_db_path)
    cur = con.cursor()
    cur.execute("CREATE TABLE frser_sqlite (timestamp INTEGER, version TEXT, downloads INTEGER)")
    con.commit()

    test_data = {
        "items": [
            {"version": "2.2.1", "downloads": 1208},
            {"version": "2.1.1", "downloads": 45478},
        ],
    }

    now_in_seconds = time.time()
    with mock.patch("main.DB_LOCATION", new=test_db_path):
        save_to_db(test_data)

    db_result = cur.execute(
        "SELECT timestamp, version, downloads INTEGER FROM frser_sqlite ORDER BY version DESC"
    ).fetchall()

    assert len(db_result) == len(test_data["items"])
    for row, item in zip(db_result, test_data["items"]):
        assert row[0] == pytest.approx(now_in_seconds, abs=1)
        assert row[1] == item["version"]
        assert row[2] == item["downloads"]
