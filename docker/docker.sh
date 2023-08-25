#!/usr/bin/env bash
## test ci
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
docker build --tag odyseeteam/commentron:$BRANCH_NAME ./
docker push odyseeteam/commentron:$BRANCH_NAME