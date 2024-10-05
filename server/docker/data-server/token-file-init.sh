#!/bin/bash

# Accept the TOKEN_PATH as an argument or environment variable
TOKEN_PATH=${1:-$TOKEN_PATH}

echo "TOKEN_PATH is: $TOKEN_PATH"

# Check if TOKEN_PATH is set
if [ -z "$TOKEN_PATH" ]; then
  echo "Error: TOKEN_PATH is not set or is empty."
  exit 1
fi

# Create the directory if it doesn't exist
mkdir -p "$(dirname "$TOKEN_PATH")"

# Attempt to create the file
if touch "$TOKEN_PATH"; then
  # Filling TOKEN_PATH with an empty valid json
  echo "{\"1\": \"363936393639\"}" >> $TOKEN_PATH

  echo "File created successfully at $TOKEN_PATH"

else
  echo "Failed to create file at $TOKEN_PATH"
  exit 1
fi
