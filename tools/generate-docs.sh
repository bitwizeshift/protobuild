#!/usr/bin/env bash

set -euo pipefail

root=$(git rev-parse --show-toplevel)
(
  cd "${root}"
  rm -rf dist htmlcov
  mkdir -p dist
  go-licenses report ./... --template "${root}/docs/licenses.md.tpl" > ${root}/docs/development/licenses.md

  mkdir -p htmlcov/html
  go test -coverprofile=htmlcov/coverage.out ./...
  go tool cover -html=htmlcov/coverage.out -o htmlcov/index.html

  mkdocs build
  GO_DOC_HTML_OUTPUT="${root}/dist/godoc" ${root}/tools/generate-godocs.sh
  mv ${root}/dist/godoc/* "${root}/dist/"
  rm -r "${root}/dist/godoc/"
)
