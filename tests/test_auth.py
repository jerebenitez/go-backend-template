import os
import bcrypt
import requests
from dotenv import load_dotenv

load_dotenv()
API_URL = (
    f"{os.getenv('HOST', default='localhost')}:{os.getenv('PORT', default='3000')}"
)


def test_signup(postgres, server, clean_db):
    salt = bcrypt.gensalt()
    password = "password123".encode("UTF-8")
    payload = {
        "email": "test@email.com",
        "password": bcrypt.hashpw(password, salt),
        "salt": salt
    }

    response = requests.post(f"{API_URL}/users", json=payload)
    assert response.status_code == 201
    assert "id" in response.json()
    assert response.json()["email"] == payload["email"]
