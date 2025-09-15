import subprocess
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


@pytest.fixture(scope="session")
def server(network, postgres):
    """Init server, connect to database"""
    port = postgres.get_exposed_port(postgres.port)
    db_url = f"postgresql://{postgres.username}:{postgres.password}@db:5432/{postgres.dbname}"

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
                    print(f"Server ready at {base_url}")
                    break
            except requests.ConnectionError as e:
                print(f"Attempt {i+1}: server not ready yet ({e})")
                print(f"Server logs: {container.get_logs()}")
                time.sleep(1)
        else:
            raise RuntimeError("Server failed to start")

        yield container, base_url


def run_migrations(dsn: str):
    """Call `migrations.py` tool to apply migrations"""
    subprocess.run(["python", "../db/migrations.py", "apply", "--dsn", dsn], check=True)


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

    dsn = connection_url.rsplit("/", 1)[0] + f"/{db_name}"
    yield dsn
