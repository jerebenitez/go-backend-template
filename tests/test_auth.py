import requests


def test_signup(server, clean_db):
    payload = {"email": "test@email.com", "password": "Password123_"}

    try:
        response = requests.post(f"{server[1]}/auth/signup", json=payload)
    except requests.ConnectionError:
        print("Server logs:\n", server[0].get_logs())
        raise

    if response.status_code == 400:
        print(response.text)
        print("Server logs:\n", server[0].get_logs())

    assert response.status_code == 201
    assert "id" in response.json()
    assert "createdAt" in response.json()
    assert response.json()["email"] == payload["email"]
