#!/bin/sh
set -eu

mkdir -p /app/logs /app/.ssh

if [ $# -eq 0 ]; then
  set -- serve
fi

if [ "$1" = "local" ]; then
  shift
fi

exec /usr/local/bin/portfolio-tui "$@"
