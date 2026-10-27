#!/usr/bin/env bash
# Read-only diagnostics for Hisense VIDAA TV
# Usage: TV_IP=192.168.1.100 ./vidaa_diag.sh

set -e
if [ -z "$TV_IP" ]; then
  echo "TV_IP required"
  exit 1
fi
echo "Checking ports on $TV_IP..."
timeout 3 bash -c "</dev/tcp/$TV_IP/36669" && echo "36669 open" || echo "36669 closed/filtered"
