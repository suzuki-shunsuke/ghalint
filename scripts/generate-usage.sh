#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/.."

command_console() {
  echo '```console'
  echo "$ $*"
  "$@"
  echo '```'
}

commands() {
  for cmd in run run-action; do
    echo "
## ghalint $cmd

$(command_console ghalint help $cmd)"
  done
}

echo -n "# Usage

<!-- This is generated by scripts/generate-usage.sh. Don't edit this file directly. -->

$(command_console ghalint help)
$(commands)
" > docs/usage.md