#!/bin/bash

set -eo pipefail

docker login -u $DOCKER_USER -p $DOCKER_PASS
curl -sL http://git.io/goreleaser | bash