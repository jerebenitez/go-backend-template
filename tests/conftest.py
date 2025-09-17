import time
import uuid

import psycopg2
import pytest
import requests
from testcontainers.core.container import DockerContainer
from testcontainers.core.network import Network
from testcontainers.postgres import PostgresContainer


@pytest.fixture(scope="session")
def network():
    with Network() as network:
        yield network


@pytest.fixture(scope="session")
def postgres(network):
    """Init Postgres 16 database"""
    with PostgresContainer("postgres:16") as pg:
        pg.with_network(network)
        pg.with_network_aliases("db")
        pg.with_name("db")

        pg.start()

        yield pg


@pytest.fixture()
def server(network, postgres, clean_db):
    """Init server, connect to database"""
    db_url = f"postgresql://{postgres.username}:{postgres.password}@db:5432/{clean_db}?sslmode=disable"

    with DockerContainer("server:latest") as container:
        container.with_network(network).with_env("DB_DSN", db_url)
        container.with_exposed_ports(5000)

        container.start()

        host = container.get_container_host_ip()
        port = container.get_exposed_port(5000)
        base_url = f"http://{host}:{port}"

        for i in range(30):
            try:
                r = requests.get(f"{base_url}/health-check")
                if r.status_code == 200:
                    print("Server ready!")
                    break
            except requests.ConnectionError as e:
                print(f"Attempt {i+1}: server not ready yet ({e})")
                print(f"Server logs: {container.get_logs()}")
                time.sleep(1)
        else:
            raise RuntimeError("Server failed to start")

        yield container, base_url


def run_migrations(dsn):
    import os
    import subprocess
    import sys
    from pathlib import Path

    print(f"Running migrations with DSN: {dsn}")
    print(f"Current working directory: {os.getcwd()}")

    # Change to the db directory where migrations.py is located
    db_dir = Path(__file__).parent.parent / "db"  # Adjust this path as needed
    migrations_script = db_dir / "migrations.py"

    if not migrations_script.exists():
        raise FileNotFoundError(f"migrations.py not found at {migrations_script}")

    # Run the migrations script from the db directory
    try:
        result = subprocess.run(
            [sys.executable, str(migrations_script), "apply", "--dsn", dsn],
            cwd=str(db_dir),
            capture_output=True,
            text=True,
        )

        print(f"Migration stdout: {result.stdout}")
        print(f"Migration stderr: {result.stderr}")
        print(f"Migration return code: {result.returncode}")

        if result.returncode != 0:
            raise RuntimeError(f"Migration failed: {result.stderr}")

    except Exception as e:
        print(f"Migration error: {e}")
        raise


@pytest.fixture(scope="session")
def initial_db(postgres):
    """Create db template with migrations and seeds"""
    template_name = "initialdb"
    connection_url = postgres.get_connection_url(driver=None)

    conn = psycopg2.connect(connection_url)
    conn.autocommit = True
    cur = conn.cursor()

    cur.execute(f"DROP DATABASE IF EXISTS {template_name}")
    cur.execute(f"CREATE DATABASE {template_name}")

    cur.close()
    conn.close()

    base_dsn = connection_url.rsplit("/", 1)[0]
    template_dsn = f"{base_dsn}/{template_name}"
    run_migrations(template_dsn)

    # Verify migrations worked
    conn = psycopg2.connect(template_dsn)
    cur = conn.cursor()
    cur.execute(
        "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
    )
    tables = cur.fetchall()
    print(f"Tables in template database: {tables}")
    cur.close()
    conn.close()

    yield template_name


@pytest.fixture
def clean_db(postgres, initial_db):
    """Clone initial_db to create a fesh db for each test"""
    connection_url = postgres.get_connection_url(driver=None)
    db_name = f"test_{uuid.uuid4().hex[:8]}"

    conn = psycopg2.connect(connection_url)
    conn.autocommit = True
    cur = conn.cursor()

    cur.execute(f'CREATE DATABASE "{db_name}" WITH TEMPLATE {initial_db}')

    cur.close()
    conn.close()

    yield db_name
