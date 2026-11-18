import uuid

import httpx

from helpers import register_user


def test_transfers_require_auth(client: httpx.Client) -> None:
    r = client.get("/api/v1/transfers")
    assert r.status_code == 401
    assert r.json()["error"] == "unauthorized"


def test_create_transfer_and_get_public(
    client: httpx.Client,
    tracked_users: list[uuid.UUID],
) -> None:
    user = register_user(client, tracked_users)
    title = "int transfer"

    create = client.post(
        "/api/v1/transfers",
        headers={"Authorization": f"Bearer {user['token']}"},
        json={"title": title, "expires_in_hours": 24},
    )
    assert create.status_code == 201, create.text
    transfer = create.json()
    assert transfer["short_code"]
    assert transfer["status"] == "active"

    pub = client.get(f"/api/v1/public/transfers/{transfer['short_code']}")
    assert pub.status_code == 200, pub.text
    body = pub.json()
    assert body["short_code"] == transfer["short_code"]
    assert body["title"] == title
    assert body["sender_email"] == user["email"]


def test_public_transfer_not_found(client: httpx.Client) -> None:
    r = client.get("/api/v1/public/transfers/doesnotexist")
    assert r.status_code == 404
