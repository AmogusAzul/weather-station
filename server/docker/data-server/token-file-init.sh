#!/bin/bash
set -e

echo "TOKEN_PATH"

# Check if TOKEN_PATH is set
if [ -z "$TOKEN_PATH" ]; then
  echo "Error: TOKEN_PATH is not set or is empty."
  exit 1
fi

# Debugging: Show the value of TOKEN_PATH
echo "TOKEN_PATH is set to: $TOKEN_PATH"

# Create the directory if it doesn't exist
mkdir -p "$(dirname "$TOKEN_PATH")"

# Attempt to create the file
if touch "$TOKEN_PATH"; then
  echo "File created successfully at $TOKEN_PATH"
else
  echo "Failed to create file at $TOKEN_PATH"
  exit 1
fi