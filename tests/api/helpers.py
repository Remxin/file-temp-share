from __future__ import annotations

import uuid

import httpx


def unique_email() -> str:
    return f"int-{uuid.uuid4()}@filetemp.local"


def register_user(
    client: httpx.Client,
    tracked_users: list[uuid.UUID],
    *,
    email: str | None = None,
    password: str = "password123",
) -> dict:
    email = email or unique_email()
    r = client.post("/api/v1/auth/register", json={"email": email, "password": password})
    assert r.status_code == 201, r.text
    body = r.json()
    tracked_users.append(uuid.UUID(body["user"]["id"]))
    return {
        "email": email,
        "password": password,
        "token": body["token"],
        "user": body["user"],
    }
