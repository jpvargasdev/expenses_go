#!/bin/sh

# Check if the environment variable is set
if [ -z "$GOOGLE_APPLICATION_CREDENTIALS_JSON" ]; then
  echo "Error: GOOGLE_APPLICATION_CREDENTIALS_JSON is not set."
  exit 1
fi

# Write the JSON string to a file
echo "$GOOGLE_APPLICATION_CREDENTIALS_JSON" | jq '.' > /app/firebase-config.json

# Export the path to make Firebase SDK recognize it
export GOOGLE_APPLICATION_CREDENTIALS="/app/firebase-config.json"

# Path to the SQLite database
DATABASE_PATH="/data/database.db"

# Ensure the volume directory exists
if [ ! -d "/data" ]; then
  echo "Creating data directory..."
  mkdir -p /data
fi

# Initialize database if it doesn't exist
if [ ! -e "$DATABASE_PATH" ]; then
  echo "Initializing database..."
  sqlite3 "$DATABASE_PATH" < init_db.sql
  echo "Database initialized."
else
  echo "Database already exists. Skipping initialization."
fi

# Seed the database if needed
SCHEMA_NOT_CHANGED_HASH="bb78723ad3f70982ed78a0da36c2ac45399aa536"
SEED_ALREADY_RAN="/data/seed_already_ran"
SCHEMA_HASH=$(sha1sum seed_db.sql | cut -d ' ' -f 1)

if [ ! -f "$SEED_ALREADY_RAN" ] && [ "$SCHEMA_NOT_CHANGED_HASH" != "$SCHEMA_HASH" ]; then
  echo "Seeding database with seed_db.sql..."
  sqlite3 "$DATABASE_PATH" < seed_db.sql
  touch "$SEED_ALREADY_RAN"
  echo "Seeding complete."
else
  echo "No changes in seed_db.sql or seed already applied."
fi

# Start the Go application and wsgi in the background
./guilliman &  # Start the Go app
python wsgi.py "$DATABASE_PATH" &  # Start sqlite_web on port 8081

# Wait for all background processes to finish
wait
