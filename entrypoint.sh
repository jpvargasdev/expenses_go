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

# Start the Go application and wsgi in the background
./guilliman  # Start the Go app

wait
