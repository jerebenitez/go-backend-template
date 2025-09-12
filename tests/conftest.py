import os
import uuid
import time
import psycopg2
import pytest
import requests
import subprocess
from testcontainers.postgres import PostgresContainer
from testcontainers.core.container import DockerContainer
from testcontainers.core.network import Network


@pytest.fixture(scope="session")
def network():
    with Network() as net:
        yield net


@pytest.fixture(scope="session")
def postgres(network):
    """Init Postgres 16 database"""
    with PostgresContainer("postgres:16") as pg:
        pg.with_network(network).with_network_aliases("db")
        pg.start()
        os.environ["DB_DSN"] = pg.get_connection_url()

        yield pg


@pytest.fixture(scope="session")
def server(postgres, network):
    """Init server, connect to database"""
    db_url = f"postgres://{postgres.username}:{postgres.password}@db:{postgres.port}/{postgres.dbname}?sslmode=disable"

    with DockerContainer("") as container:
        container.with_network(network) \
                .with_env("DB_DSN", db_url) \
                .with_exposed_ports(5000) \
                .start()
        host = container.get_container_host_ip()
        port = container.get_exposed_port(5000)
        base_url = f"http://{host}:{port}"

        for _ in range(30):
            try:
                requests.get(f"{base_url}/health-check")
                break
            except Exception:
                time.Sleep(1)
        else:
            raise RuntimeError("Server failed to start")

        yield base_url


def run_migrations(dsn: str):
    """Call `migrations.py` tool to apply migrations"""
    subprocess.run(
        ["python", "../migrations.py", "apply", "--dsn", dsn],
        check=True
    )


@pytest.fixture(scope="session")
def initial_db(postgres):
    """Create db template with migrations and seeds"""
    template_name = "initialDb"

    conn = psycopg2.connect(postgres.get_connection_url())
    conn.autocommit = True
    cur = conn.cursor()

    cur.execute(f"DROP DATABASE IF EXISTS {template_name}")
    cur.execute(f"CREATE DATABASE {template_name}")

    cur.close()
    conn.close()

    base_dsn = postgres.get_connection_url().split("/", 1)[0]
    template_dsn = f"{base_dsn}/{template_name}"
    run_migrations(template_dsn)

    yield template_name


@pytest.fixture
def clean_db(postgres, initial_db):
    """Clone initial_db to create a fesh db for each test"""
    db_name = f"test_{uuid.uuid4().hex[:8]}"
    conn = psycopg2.connect(postgres.get_connection_url())
    conn.autocommit = True
    cur = conn.cursor()

    cur.execute(f'CREATE DATABASE "{db_name} WITH TEMPLATE {initial_db}')

    cur.close()
    conn.close()

    dsn = postgres.get_connection_url().rsplit("/", 1)[0] + f"/{db_name}"
    yield dsn
