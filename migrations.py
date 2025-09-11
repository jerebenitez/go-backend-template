import os
import sys
import psycopg2
from pathlib import Path
from dotenv import load_dotenv


def get_template(idx: int) -> str:
    template = f"""BEGIN;
INSERT INTO migrations(ref) VALUES ({idx});

-- Insert migration code here

COMMIT;
"""
    return template


def print_usage():
    print(
        f"""Usage: {sys.argv[0]} CMD [OPTIONS]

Commands:
apply\t\tApply existing migrations
new [name] \tCreate new migration file with 'name' as prefix"""
    )


def get_next_migration_number(dir: Path) -> str:
    last_migration = 0
    for f in dir.glob(".sql"):
        name = f.stem
        try:
            num = int(name.split("_")[-1])
            last_migration = num if num > last_migration else last_migration
        except ValueError:
            # Skip files that don't follow the name_####.sql format
            continue

    return f"{last_migration+1:04d}"


def create_new_migration(name: str, dir: str):
    migrations_dir = Path(dir)
    migrations_dir.mkdir(parents=True, exist_ok=True)

    idx = get_next_migration_number(migrations_dir)

    migration_file = migrations_dir / f"{name}_{idx}.sql"
    if migration_file.exists():
        print(f"Error: migration file {migration_file} already exists!")
        sys.exit(1)

    migration_file.write_text(get_template(int(idx)))
    print(f"Created migration: {migration_file}")


def apply_migrations(dsn: str, dir: str):
    migrations_dir = Path(dir)

    conn = psycopg2.connect(dsn)
    conn.autocommit = False
    cur = conn.cursor()

    # Ensure migrations table exists, used to track which migrations
    # have been applied
    cur.execute("""
        CREATE TABLE IF NOT EXISTS migrations (
            id SERIAL KEY,
            ref INTEGER,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    """)
    conn.commit()

    # Get applied migrations
    cur.execute("SELECT ref FROM migrations")
    applied = {row[0] for row in cur.fetchall()}

    migration_files = sorted(
        migrations_dir.glob("*.sql"),
        key=lambda f: int(f.stem.split("_")[-1])
    )

    for f in migration_files:
        idx = int(f.stem.split("_")[-1])
        if idx in applied:
            continue

        print(f"Applying migration {f.name}")
        sql = f.read_text()
        try:
            cur.execute(sql)
            conn.commit()
        except Exception as e:
            conn.rollback()
            print(f"Failed to apply {f.name}: {e}")

    cur.close()
    conn.close()


def get_dsn_from_env(dir=".env") -> str:
    """
    Expected keys: DB_NAME, DB_USER, DB_PASSWORD, DB_HOST, DB_PORT
    """
    env_file = Path(dir)
    load_dotenv(dotenv_path=env_file)

    dsn = (
        f"dbname={os.getenv('DB_NAME', 'postgres')} "
        f"user={os.getenv('DB_USER', 'postgres')} "
        f"password={os.getenv('DB_PASSWORD')} "
        f"host={os.getenv('DB_HOST', 'localhost')} "
        f"port={os.getenv('DB_PORT', '5432')}"
    )

    return dsn


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print_usage()
        sys.exit(1)

    if sys.argv[1] == "apply":
        apply_migrations(get_dsn_from_env(), "db/migrations")
    elif sys.argv[1] == "new" and len(sys.argv) == 3:
        create_new_migration(sys.argv[2], "db/migrations/")
    else:
        print_usage()
        sys.exit(1)
