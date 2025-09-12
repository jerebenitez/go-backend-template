import bcrypt
import requests


def test_signup(server, clean_db):
    salt = bcrypt.gensalt()
    password = "password123".encode("UTF-8")
    hased_pwd = bcrypt.hashpw(password, salt)

    payload = {
        "email": "test@email.com",
        "password": hased_pwd.decode("utf-8", errors="ignore"),
        "salt": salt.decode("utf-8", errors="ignore"),
    }

    try:
        response = requests.post(f"{server[1]}/users", json=payload)
    except requests.ConnectionError:
        print("Server logs:\n", server[0].get_logs())
        raise
    assert response.status_code == 201
    assert "id" in response.json()
    assert response.json()["email"] == payload["email"]
