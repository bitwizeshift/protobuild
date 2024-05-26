#!/usr/bin/env bash

set -euo pipefail

root=$(git rev-parse --show-toplevel)
(
  cd "${root}"
  mkdir -p dist
  go-licenses report ./... --template "${root}/docs/licenses.md.tpl" > ${root}/docs/development/licenses.md

  mkdocs build
)
