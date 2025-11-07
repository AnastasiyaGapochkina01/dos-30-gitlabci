#!/bin/bash
echo "Run application on port 9100"
nohup python3 -m http.server 9100 >/dev/null &
