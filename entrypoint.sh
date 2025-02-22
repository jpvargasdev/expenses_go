#!/bin/sh

# Check if the environment variable is set
if [ -z "$GOOGLE_APPLICATION_CREDENTIALS_JSON" ]; then
  echo "Error: GOOGLE_APPLICATION_CREDENTIALS_JSON is not set."
  exit 1
fi

# Decode the Base64 string and write it to a file
echo "$GOOGLE_APPLICATION_CREDENTIALS_JSON" | base64 -d > /app/firebase-config.json

# Export the path to make Firebase SDK recognize it
export GOOGLE_APPLICATION_CREDENTIALS="/app/firebase-config.json"

# Start the Go application
./guilliman  # Start the Go app

wait
