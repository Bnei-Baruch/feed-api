#!/usr/bin/env bash
# run misc/build.sh from project root
set -e
set -x

docker image build -t feed-api:latest .
version="$(docker run --rm feed-api:latest /app/feed-api version | awk '{print $NF}')"
docker create --name dummy feed-api:latest
docker cp dummy:/app/feed-api ./feed-api-linux-"${version}"
docker rm -f dummy
