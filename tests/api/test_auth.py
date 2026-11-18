import uuid

import httpx

from helpers import register_user


def test_register_login_and_me(
    client: httpx.Client,
    tracked_users: list[uuid.UUID],
) -> None:
    user = register_user(client, tracked_users)
    email = user["email"]
    password = user["password"]

    login = client.post(
        "/api/v1/auth/login",
        json={"email": email, "password": password},
    )
    assert login.status_code == 200, login.text
    token = login.json()["token"]
    assert token

    me = client.get("/api/v1/auth/me", headers={"Authorization": f"Bearer {token}"})
    assert me.status_code == 200, me.text
    body = me.json()
    assert body["email"] == email
    assert body["is_active"] is True


def test_login_rejects_bad_password(
    client: httpx.Client,
    tracked_users: list[uuid.UUID],
) -> None:
    user = register_user(client, tracked_users)
    r = client.post(
        "/api/v1/auth/login",
        json={"email": user["email"], "password": "wrong-password"},
    )
    assert r.status_code == 401


def test_register_conflict_on_duplicate_email(
    client: httpx.Client,
    tracked_users: list[uuid.UUID],
) -> None:
    user = register_user(client, tracked_users)
    r = client.post(
        "/api/v1/auth/register",
        json={"email": user["email"], "password": user["password"]},
    )
    assert r.status_code == 409


def test_change_password_then_login(
    client: httpx.Client,
    tracked_users: list[uuid.UUID],
) -> None:
    user = register_user(client, tracked_users)
    new_password = "new-password-99"

    change = client.post(
        "/api/v1/auth/password",
        headers={"Authorization": f"Bearer {user['token']}"},
        json={
            "current_password": user["password"],
            "new_password": new_password,
        },
    )
    assert change.status_code == 204, change.text

    old_login = client.post(
        "/api/v1/auth/login",
        json={"email": user["email"], "password": user["password"]},
    )
    assert old_login.status_code == 401

    new_login = client.post(
        "/api/v1/auth/login",
        json={"email": user["email"], "password": new_password},
    )
    assert new_login.status_code == 200, new_login.text
