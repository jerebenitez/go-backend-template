import sys
from pathlib import Path


def get_template(idx: int) -> str:
    template = f"""BEGIN;
INSERT INTO migrations(ref) VALUES ({idx});

-- Insert migration code here

COMMIT;
"""
    return template


def print_usage():
    print(f"""Usage: {sys.argv[0]} CMD [OPTIONS]

Commands:
apply\t\tApply existing migrations
new [name] \tCreate new migration file with 'name' as prefix""")


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


def create_new_migration(name: str):
    migrations_dir = Path("db/migrations/")
    idx = get_next_migration_number(migrations_dir)

    migration_file = migrations_dir / f"{name}_{idx}.sql"
    if migration_file.exists():
        print(f"Error: migration file {migration_file} already exists!")
        sys.exit(1)

    migration_file.write_text(get_template(int(idx)))
    print(f"Created migration: {migration_file}")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print_usage()
        sys.exit(1)

    if sys.argv[1] == "apply":
        pass
    elif sys.argv[1] == "new" and len(sys.argv) == 3:
        create_new_migration(sys.argv[2])
    else:
        print_usage()
        sys.exit(1)
