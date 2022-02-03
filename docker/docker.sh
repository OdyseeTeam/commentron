#!/usr/bin/env bash
## test ci
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
docker build --tag odyseeteam/commentron:$TRAVIS_BRANCH ./
docker push odyseeteam/commentron:$TRAVIS_BRANCH