#!/usr/bin/env bash

err=0
trap 'err=1' ERR
# All the .go files, excluding vendor/ and model (auto generated)
GO_FILES=$(find . -iname '*.go' -type f | grep -v /model/*  | grep -v /migration/*)
(
	go install golang.org/x/tools/cmd/goimports@latest                   # Used in build script for generated files
	go install golang.org/x/lint/golint@latest                           # Linter
	go install github.com/mdempsky/unconvert@latest                      # Identifies unnecessary type conversions
	go install github.com/kisielk/errcheck@latest                        # Checks for unhandled errors
)
# go vet is the official Go static analyzer
echo "Running go vet..." && go vet $(go list ./... | grep -v /migration/* | grep -v /model/* | grep -v /server/services/v1/rpc/*  )
# checks for unhandled errors
echo "Running errcheck..." && errcheck $(go list ./... | grep -v /migration/* | grep -v /model/*  | grep -v /server/services/v1/rpc/* )
# check for unnecessary conversions - ignore autogen code
echo "Running unconvert..." && unconvert -v $(go list ./... | grep -v /migration/* | grep -v /model/* | grep -v /server/services/v1/rpc/*  )
# one last linter - ignore autogen code
echo "Running golint..." && golint -set_exit_status $(go list ./... | grep -v /migration/* | grep -v /model/* | grep -v /server/services/v1/rpc/*  )
test $err = 0 # Return non-zero if any command failed