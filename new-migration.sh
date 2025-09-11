##!/usr/bin/env bash
set -euo pipefail

# Check if first argument is provided
if [[ -z "${1:-}" ]]; then
    echo "Usage: $0 <migration_name>"
    exit 1
fi

migration_name="$1"
last_migration=0

shopt -s nullglob
for f in db/migrations/*.sql; do
    num="${f##*_}"
    num="${num%.sql}"
    (( num > last_migration )) && last_migration=$num
done
shopt -u nullglob

last_migration=$((10#$last_migration+1))
last_migration=$(printf "%04d" "$last_migration")

migration_file="${migration_name}_${last_migration}.sql"

cat <<EOF > "db/migrations/$migration_file"
BEGIN;
INSERT INTO migrations (ref) VALUES (1);

-- Insert migration code here

COMMIT;
EOF

echo "Created migration: db/migrations/$migration_file"
