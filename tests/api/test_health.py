import httpx


def test_health_ok(client: httpx.Client) -> None:
    r = client.get("/health")
    assert r.status_code == 200
    assert r.json()["status"] == "ok"
