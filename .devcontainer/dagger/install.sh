#!/usr/bin/env bash

DAGGER_VERSION="v${VERSION:-"0.3.10"}"

cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | sh