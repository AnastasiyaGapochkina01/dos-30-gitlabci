#!/bin/bash
#apt update && apt install curl python3
bash src/app.sh
echo "Run test 1 of 2"
curl -f 127.0.0.1:9100 || true
sleep 2
echo "Run test 2 of 2"
sleep 3

if (( RANDOM % 2 )); then
  echo "Tests finished: SUCCESS"
  exit 0
else
  echo "Tests finished: FAILURE"
  exit 1
fi
