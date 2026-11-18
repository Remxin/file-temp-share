from __future__ import annotations

import os
import time
import uuid
from collections.abc import Iterator

import httpx
import psycopg
import pytest

API_BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8080").rstrip("/")
DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgres://filetemp:filetemp@localhost:5432/filetemp?sslmode=disable",
)


@pytest.fixture(scope="session")
def api_base_url() -> str:
    return API_BASE_URL


@pytest.fixture(scope="session")
def wait_for_api(api_base_url: str) -> None:
    deadline = time.time() + 60
    last_err: Exception | None = None
    while time.time() < deadline:
        try:
            with httpx.Client(base_url=api_base_url, timeout=2.0) as client:
                r = client.get("/health")
                if r.status_code == 200:
                    return
        except Exception as exc:  # noqa: BLE001 — retry until ready
            last_err = exc
        time.sleep(1)
    raise RuntimeError(f"API not ready at {api_base_url}: {last_err}")


@pytest.fixture(scope="session")
def db() -> Iterator[psycopg.Connection]:
    conn = psycopg.connect(DATABASE_URL, autocommit=True)
    try:
        yield conn
    finally:
        conn.close()


@pytest.fixture
def client(api_base_url: str, wait_for_api: None) -> Iterator[httpx.Client]:
    with httpx.Client(base_url=api_base_url, timeout=10.0) as http:
        yield http


@pytest.fixture
def tracked_users(db: psycopg.Connection) -> Iterator[list[uuid.UUID]]:
    """Collect user IDs created during a test and delete them on teardown (CASCADE)."""
    ids: list[uuid.UUID] = []
    yield ids
    for user_id in ids:
        db.execute("DELETE FROM users WHERE id = %s", (user_id,))
