#!/bin/sh

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <json_file>"
  exit 1
fi

jq -c '.' "$1"
