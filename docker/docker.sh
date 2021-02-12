#!/usr/bin/env bash

echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
docker build --tag lbry/commentron:$TRAVIS_BRANCH ./
docker push lbry/commentron:$TRAVIS_BRANCH