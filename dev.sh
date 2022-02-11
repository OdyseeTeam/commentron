#!/bin/bash
set -euo pipefail
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_DIR="$DIR"

(
  source "$DIR/scripts/setup.sh"
  touch -a .env && set -o allexport
  source ./.env; set +o allexport

  hash reflex 2>/dev/null || go get github.com/cespare/reflex
  hash reflex 2>/dev/null || { echo >&2 'Make sure $GOPATH/bin is in your $PATH'; exit 1;  }

  hash go-bindata 2>/dev/null || go install github.com/jteeuwen/go-bindata/...@latest

  cd "$APP_DIR"
  #golint -set_exit_status $(go list ./... | grep -v /migration/* )
  reflex --decoration=none --start-service=true --regex='\.go$' --inverse-regex='migration/bindata\.go' -- sh -c "go generate && go run *.go serve -d"
)