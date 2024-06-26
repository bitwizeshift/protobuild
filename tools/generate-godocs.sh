#!/usr/bin/env bash

# Source: https://gitlab.com/-/snippets/1977141

set -euo pipefail

function extract_module_name {
  # Extract module name
  sed -n -E 's/^\s*module\s+([[:graph:]]+)\s*$/\1/p'
}

function normalize_url {
  # Normalize provided URL. Removing double slashes
  echo "$1" | sed -E 's,([^:]/)/+,\1,g'
}

GO_MODULE="github.com/bitwizeshift/protobuild"

function generate_go_documentation {
  # Go doc
  local URL
  local PID
  local STATUS

  # Setup
  rm -rf "${GO_DOC_HTML_OUTPUT:-godoc}"
  mkdir -p "${GO_DOC_HTML_OUTPUT:-godoc}"

  # URL path to Go package and module documentation
  URL=$(normalize_url "http://${GO_DOC_HTTP:-localhost:6060}/pkg/$GO_MODULE/")

  # Starting godoc server
  echo "Starting godoc server..."
  godoc -http="${GO_DOC_HTTP:-localhost:6060}" &
  PID=$!

  # Waiting for godoc server
  while ! curl --fail --silent "$URL" 2>&1 >/dev/null; do
    sleep 0.1
  done

  # Download all documentation content from running godoc server
  wget \
    --recursive \
    --no-verbose \
    --convert-links \
    --page-requisites \
    --adjust-extension \
    --execute=robots=off \
    --include-directories="/lib,/pkg/$GO_MODULE,/src/$GO_MODULE" \
    --exclude-directories="*" \
    --directory-prefix="${GO_DOC_HTML_OUTPUT:-godoc}" \
    --no-host-directories \
    "$URL"

  # Stop godoc server
  kill -9 "$PID"
  echo "Stopped godoc server"
  echo "Go source code documentation generated under ${GO_DOC_HTML_OUTPUT:-godoc}"
  echo "Package docs available at '$(pwd)/${GO_DOC_HTML_OUTPUT:-godoc}/pkg/${GO_MODULE}'"
}
generate_go_documentation
